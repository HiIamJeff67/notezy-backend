import { createServer } from "node:http";

import { WebSocket, WebSocketServer } from "ws";
import * as Y from "yjs";

import { BlockNoteProjector } from "./blocknote_projector.js";
import { config } from "./config.js";
import  { type InternalFrame, createInternalFrame, parseInternalFrame } from "./types/internal_frame.js";
import { InternalFrameType } from "./types/internal_frame_type.js";
import type { Room } from "./types/room.js";
import { parseYjsDocumentState, parseYjsUpdateSequence } from "./types/yjs_document_state.js";
import { RoomRegistry } from "./room_registry.js";

const roomRegistry = new RoomRegistry();
const webSocketServer = new WebSocketServer({ noServer: true });
const blockNoteProjector = new BlockNoteProjector();
const projectionDebounceMilliseconds = 300;
const projectionRetryMilliseconds = 1_000;

/* ============================== Internal frame delivery ============================== */

function broadcastInternalFrame(
  room: Room,
  type: InternalFrameType,
  blockPackId: string,
  payload: Buffer,
): void {
  for (const subscriber of room.subscribers.values()) {
    sendInternalFrame(
      subscriber.webSocket,
      type,
      subscriber.connectionId,
      subscriber.connectorChannelId,
      blockPackId,
      payload,
    );
  }
}

function sendInternalFrame(
  webSocket: WebSocket,
  type: InternalFrameType,
  connectionId: string,
  connectorChannelId: number,
  blockPackId: string,
  payload: Buffer = Buffer.alloc(0),
): boolean {
  if (webSocket.readyState !== WebSocket.OPEN) {
    return false;
  }

  webSocket.send(
    createInternalFrame(type, connectionId, connectorChannelId, blockPackId, payload),
  );

  return true;
}

function sendRoomInitialState(room: Room, blockPackId: string): void {
  if (room.document === null) {
    return;
  }

  broadcastInternalFrame(
    room,
    InternalFrameType.InternalFrameType_YjsDocument,
    blockPackId,
    Buffer.from(Y.encodeStateAsUpdate(room.document)),
  );
}

function resyncRoom(room: Room, blockPackId: string): void {
  if (room.projectionTimer !== null) {
    clearTimeout(room.projectionTimer);
  }

  room.document = null;
  room.isLoading = false;
  room.dirtyUpdateCount = 0;
  room.lastUpdateSequence = 0;
  room.compactedUntilSequence = 0;
  room.projectedUntilSequence = -1;
  room.pendingYjsUpdates = [];
  room.inFlightYjsUpdate = null;
  room.projectionTimer = null;
  room.inFlightProjection = null;

  broadcastInternalFrame(
    room,
    InternalFrameType.InternalFrameType_ResyncRequired,
    blockPackId,
    Buffer.alloc(0),
  );
}

function requestRoomLoad(room: Room, frame: InternalFrame, webSocket: WebSocket): void {
  if (room.isLoading) {
    return;
  }

  room.isLoading = true;
  if (sendInternalFrame(
    webSocket,
    InternalFrameType.InternalFrameType_LoadYjsDocument,
    frame.connectionId,
    frame.connectorChannelId,
    frame.blockPackId,
  )) {
    return;
  }

  resyncRoom(room, frame.blockPackId);
}

function processNextYjsUpdate(room: Room): void {
  if (room.document === null || room.inFlightYjsUpdate !== null) {
    return;
  }

  const pendingYjsUpdate = room.pendingYjsUpdates.shift();
  if (pendingYjsUpdate === undefined) {
    return;
  }

  try {
    Y.applyUpdate(room.document, pendingYjsUpdate.frame.payload);
  } catch {
    sendInternalFrame(
      pendingYjsUpdate.webSocket,
      InternalFrameType.InternalFrameType_ResyncRequired,
      pendingYjsUpdate.frame.connectionId,
      pendingYjsUpdate.frame.connectorChannelId,
      pendingYjsUpdate.frame.blockPackId,
    );

    processNextYjsUpdate(room);

    return;
  }

  room.inFlightYjsUpdate = pendingYjsUpdate;
  if (sendInternalFrame(
    pendingYjsUpdate.webSocket,
    InternalFrameType.InternalFrameType_AppendYjsUpdate,
    pendingYjsUpdate.frame.connectionId,
    pendingYjsUpdate.frame.connectorChannelId,
    pendingYjsUpdate.frame.blockPackId,
    pendingYjsUpdate.frame.payload,
  )) {
    return;
  }

  resyncRoom(room, pendingYjsUpdate.frame.blockPackId);
}

function scheduleBlockProjection(
  room: Room,
  blockPackId: string,
  delayMilliseconds: number = projectionDebounceMilliseconds,
): void {
  if (
    room.document === null ||
    room.inFlightYjsUpdate !== null ||
    room.pendingYjsUpdates.length > 0 ||
    room.inFlightProjection !== null ||
    room.projectionTimer !== null ||
    room.lastUpdateSequence <= room.projectedUntilSequence ||
    room.subscribers.size === 0
  ) {
    return;
  }

  room.projectionTimer = setTimeout(() => {
    room.projectionTimer = null;

    if (
      room.document === null ||
      room.inFlightYjsUpdate !== null ||
      room.pendingYjsUpdates.length > 0 ||
      room.inFlightProjection !== null ||
      room.lastUpdateSequence <= room.projectedUntilSequence
    ) {
      return;
    }

    const subscriber = room.subscribers.values().next().value;
    if (subscriber === undefined) {
      return;
    }

    const projectedSequence = room.lastUpdateSequence;
    let payload: Buffer;
    try {
      payload = Buffer.from(JSON.stringify({
        schemaId: "notezy.blocknote",
        schemaVersion: 1,
        projectedSequence,
        blocks: blockNoteProjector.projectYjsDocument(room.document),
      }));
    } catch (error) {
      console.error("failed to project Yjs document", {
        blockPackId,
        error,
      });
      scheduleBlockProjection(room, blockPackId, projectionRetryMilliseconds);

      return;
    }

    room.inFlightProjection = {
      connectionId: subscriber.connectionId,
      connectorChannelId: subscriber.connectorChannelId,
      projectedSequence,
    };
    if (sendInternalFrame(
      subscriber.webSocket,
      InternalFrameType.InternalFrameType_ApplyBlockProjection,
      subscriber.connectionId,
      subscriber.connectorChannelId,
      blockPackId,
      payload,
    )) {
      return;
    }

    room.inFlightProjection = null;
    scheduleBlockProjection(room, blockPackId, projectionRetryMilliseconds);
  }, delayMilliseconds);
}

function parseProjectedUntilSequence(payload: Buffer): number | null {
  try {
    const value: unknown = JSON.parse(payload.toString("utf8"));
    if (
      value === null ||
      typeof value !== "object" ||
      !("projectedUntilSequence" in value) ||
      typeof value.projectedUntilSequence !== "number" ||
      !Number.isSafeInteger(value.projectedUntilSequence) ||
      value.projectedUntilSequence < -1
    ) {
      return null;
    }

    return value.projectedUntilSequence;
  } catch {
    return null;
  }
}

/* ============================== HTTP server ============================== */

const server = createServer((request, response) => {
  const requestUrl = new URL(request.url ?? "/", `http://${request.headers.host ?? "localhost"}`);
  if (request.method === "GET" && requestUrl.pathname === "/healthz") {
    response.writeHead(200, { "content-type": "application/json" });
    response.end(JSON.stringify({ status: "ok", activeRoomCount: roomRegistry.size }));

    return;
  }

  response.writeHead(404);
  response.end();
});

/* ============================== WebSocket upgrade ============================== */

server.on("upgrade", (request, socket, head) => {
  const requestUrl = new URL(request.url ?? "/", `http://${request.headers.host ?? "localhost"}`);
  if (requestUrl.pathname !== "/internal/realtime/v1") {
    socket.destroy();

    return;
  }

  webSocketServer.handleUpgrade(request, socket, head, (webSocket) => {
    webSocketServer.emit("connection", webSocket, request);
  });
});

/* ============================== WebSocket connection ============================== */

webSocketServer.on("connection", (webSocket) => {
  webSocket.on("close", () => {
    for (const { blockPackId, room } of roomRegistry.detachAll(webSocket)) {
      resyncRoom(room, blockPackId);
    }
  });

  webSocket.on("message", (payload, isBinary) => {
    if (!isBinary) {
      webSocket.close(1003, "internal realtime frames must be binary");

      return;
    }

    let framePayload: Buffer;
    if (Buffer.isBuffer(payload)) {
      framePayload = payload;
    } else if (payload instanceof ArrayBuffer) {
      framePayload = Buffer.from(payload);
    } else if (Array.isArray(payload)) {
      framePayload = Buffer.concat(payload);
    } else {
      webSocket.close(1002, "invalid internal realtime frame");

      return;
    }

    const frame = parseInternalFrame(framePayload);
    if (frame === null || frame.version !== 1) {
      webSocket.close(1002, "invalid internal realtime frame");

      return;
    }

    switch (frame.type) {
      case InternalFrameType.InternalFrameType_Attach: {
        const room = roomRegistry.attach(
          frame.blockPackId,
          webSocket,
          frame.connectionId,
          frame.connectorChannelId,
        );
        if (room.document !== null) {
          sendInternalFrame(
            webSocket,
            InternalFrameType.InternalFrameType_YjsDocument,
            frame.connectionId,
            frame.connectorChannelId,
            frame.blockPackId,
            Buffer.from(Y.encodeStateAsUpdate(room.document)),
          );

          return;
        }

        requestRoomLoad(room, frame, webSocket);

        return;
      }
      case InternalFrameType.InternalFrameType_Detach:
        roomRegistry.detach(frame.blockPackId, frame.connectionId, frame.connectorChannelId);

        return;
      case InternalFrameType.InternalFrameType_YjsDocument: {
        const room = roomRegistry.getSubscriber(
          frame.blockPackId,
          frame.connectionId,
          frame.connectorChannelId,
        );
        if (room === undefined) {
          sendInternalFrame(
            webSocket,
            InternalFrameType.InternalFrameType_ResyncRequired,
            frame.connectionId,
            frame.connectorChannelId,
            frame.blockPackId,
          );

          return;
        }

        room.pendingYjsUpdates.push({ webSocket, frame });
        processNextYjsUpdate(room);

        return;
      }
      case InternalFrameType.InternalFrameType_Awareness: {
        const room = roomRegistry.getSubscriber(
          frame.blockPackId,
          frame.connectionId,
          frame.connectorChannelId,
        );
        if (room === undefined) {
          sendInternalFrame(
            webSocket,
            InternalFrameType.InternalFrameType_ResyncRequired,
            frame.connectionId,
            frame.connectorChannelId,
            frame.blockPackId,
          );

          return;
        }

        broadcastInternalFrame(room, frame.type, frame.blockPackId, frame.payload);

        return;
      }
      case InternalFrameType.InternalFrameType_YjsDocumentLoaded: {
        const room = roomRegistry.get(frame.blockPackId);
        const state = parseYjsDocumentState(frame.payload);
        if (room === undefined || state === null) {
          if (room !== undefined) {
            resyncRoom(room, frame.blockPackId);
          }

          return;
        }

        try {
          const document = new Y.Doc();
          if (state.snapshot.length > 0) {
            Y.applyUpdate(document, state.snapshot);
          }
          for (const update of state.updates) {
            Y.applyUpdate(document, update.payload);
          }

          room.document = document;
          room.isLoading = false;
          room.lastUpdateSequence = state.lastUpdateSequence;
          room.compactedUntilSequence = state.compactedUntilSequence;
          room.projectedUntilSequence = state.projectedUntilSequence;
          sendRoomInitialState(room, frame.blockPackId);
          processNextYjsUpdate(room);
          scheduleBlockProjection(room, frame.blockPackId);
        } catch {
          resyncRoom(room, frame.blockPackId);
        }

        return;
      }
      case InternalFrameType.InternalFrameType_YjsUpdatePersisted: {
        const room = roomRegistry.get(frame.blockPackId);
        const updateSequence = parseYjsUpdateSequence(frame.payload);
        if (
          room === undefined ||
          updateSequence === null ||
          room.inFlightYjsUpdate === null ||
          room.inFlightYjsUpdate.frame.connectionId !== frame.connectionId ||
          room.inFlightYjsUpdate.frame.connectorChannelId !== frame.connectorChannelId ||
          updateSequence !== room.lastUpdateSequence + 1
        ) {
          if (room !== undefined) {
            resyncRoom(room, frame.blockPackId);
          }

          return;
        }

        const inFlightYjsUpdate = room.inFlightYjsUpdate;
        room.inFlightYjsUpdate = null;
        room.lastUpdateSequence = updateSequence;
        room.dirtyUpdateCount += 1;
        broadcastInternalFrame(
          room,
          InternalFrameType.InternalFrameType_YjsDocument,
          frame.blockPackId,
          inFlightYjsUpdate.frame.payload,
        );
        processNextYjsUpdate(room);
        scheduleBlockProjection(room, frame.blockPackId);

        return;
      }
      case InternalFrameType.InternalFrameType_YjsPersistenceFailed: {
        const room = roomRegistry.get(frame.blockPackId);
        if (room !== undefined) {
          resyncRoom(room, frame.blockPackId);
        }

        return;
      }
      case InternalFrameType.InternalFrameType_BlockProjectionApplied: {
        const room = roomRegistry.get(frame.blockPackId);
        const projectedUntilSequence = parseProjectedUntilSequence(frame.payload);
        if (
          room === undefined ||
          projectedUntilSequence === null ||
          room.inFlightProjection === null ||
          room.inFlightProjection.connectionId !== frame.connectionId ||
          room.inFlightProjection.connectorChannelId !== frame.connectorChannelId ||
          projectedUntilSequence < room.inFlightProjection.projectedSequence ||
          projectedUntilSequence > room.lastUpdateSequence
        ) {
          if (room !== undefined) {
            resyncRoom(room, frame.blockPackId);
          }

          return;
        }

        room.inFlightProjection = null;
        room.projectedUntilSequence = projectedUntilSequence;
        scheduleBlockProjection(room, frame.blockPackId);

        return;
      }
      case InternalFrameType.InternalFrameType_BlockProjectionFailed: {
        const room = roomRegistry.get(frame.blockPackId);
        if (
          room === undefined ||
          room.inFlightProjection === null ||
          room.inFlightProjection.connectionId !== frame.connectionId ||
          room.inFlightProjection.connectorChannelId !== frame.connectorChannelId
        ) {
          return;
        }

        room.inFlightProjection = null;
        scheduleBlockProjection(room, frame.blockPackId, projectionRetryMilliseconds);

        return;
      }
      default:
        console.warn("received internal frame before its handler is enabled", {
          type: frame.type,
          blockPackId: frame.blockPackId,
        });
    }
  });
});

/* ============================== Process lifecycle ============================== */

server.listen(config.port, config.host, () => {
  console.info(`yjs worker listening on ${config.host}:${config.port}`);
});

function shutdown(signal: string): void {
  console.info(`received ${signal}, stopping yjs worker`);
  webSocketServer.clients.forEach((webSocket) => {
    webSocket.close(1001, "server shutdown");
  });

  server.close(() => process.exit(0));
}

process.once("SIGINT", () => shutdown("SIGINT"));
process.once("SIGTERM", () => shutdown("SIGTERM"));
