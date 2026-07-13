import { createServer } from "node:http";

import { WebSocket, WebSocketServer } from "ws";
import * as Y from "yjs";

import { BatchYjsHandler } from "./batch_yjs_handler.js";
import { BlockNoteProjector } from "./blocknote_projector.js";
import { config } from "./config.js";
import { YjsPersistenceBatchShutdownTimeoutMilliseconds } from "./constants/yjs_persistence_batch.js";
import { RoomRegistry } from "./room_registry.js";
import {
  createInternalFrame,
  type InternalFrame,
  parseInternalFrame,
} from "./types/internal_frame.js";
import { InternalFrameType } from "./types/internal_frame_type.js";
import type { Room } from "./types/room.js";
import {
  parseYjsDocumentState,
  parseYjsUpdateSequence,
} from "./types/yjs_document_state.js";

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
  payload: Buffer
): void {
  for (const subscriber of room.subscribers.values()) {
    if (
      !subscriber.isReady &&
      (type === InternalFrameType.InternalFrameType_YjsDocument ||
        type === InternalFrameType.InternalFrameType_Awareness)
    ) {
      continue;
    }

    sendInternalFrame(
      subscriber.webSocket,
      type,
      subscriber.connectionId,
      subscriber.connectorChannelId,
      blockPackId,
      payload
    );
  }
}

function sendInternalFrame(
  webSocket: WebSocket,
  type: InternalFrameType,
  connectionId: string,
  connectorChannelId: number,
  blockPackId: string,
  payload: Buffer = Buffer.alloc(0)
): boolean {
  if (webSocket.readyState !== WebSocket.OPEN) {
    return false;
  }

  webSocket.send(
    createInternalFrame(
      type,
      connectionId,
      connectorChannelId,
      blockPackId,
      payload
    )
  );

  return true;
}

function sendRoomInitialState(room: Room, blockPackId: string): void {
  if (room.document === null) {
    return;
  }

  const payload = Buffer.from(Y.encodeStateAsUpdate(room.document));
  for (const subscriber of room.subscribers.values()) {
    if (subscriber.isReady) {
      continue;
    }

    if (
      sendInternalFrame(
        subscriber.webSocket,
        InternalFrameType.InternalFrameType_YjsDocument,
        subscriber.connectionId,
        subscriber.connectorChannelId,
        blockPackId,
        payload
      )
    ) {
      subscriber.isReady = true;
    }
  }
}

function resyncRoom(room: Room, blockPackId: string): void {
  if (room.projectionTimer !== null) {
    clearTimeout(room.projectionTimer);
  }
  if (room.persistenceDebounceTimer !== null) {
    clearTimeout(room.persistenceDebounceTimer);
  }
  if (room.persistenceMaximumWaitTimer !== null) {
    clearTimeout(room.persistenceMaximumWaitTimer);
  }
  if (room.persistenceRetryTimer !== null) {
    clearTimeout(room.persistenceRetryTimer);
  }

  room.document = null;
  room.isLoading = false;
  room.dirtyUpdateCount = 0;
  room.lastUpdateSequence = 0;
  room.compactedUntilSequence = 0;
  room.projectedUntilSequence = -1;
  room.pendingYjsUpdates = [];
  room.pendingPersistenceUpdates = [];
  room.pendingPersistencePayloadBytes = 0;
  room.persistenceDebounceTimer = null;
  room.persistenceMaximumWaitTimer = null;
  room.persistenceRetryTimer = null;
  room.inFlightPersistenceBatch = null;
  room.projectionTimer = null;
  room.inFlightProjection = null;

  broadcastInternalFrame(
    room,
    InternalFrameType.InternalFrameType_ResyncRequired,
    blockPackId,
    Buffer.alloc(0)
  );
}

function scheduleBlockProjection(
  room: Room,
  blockPackId: string,
  delayMilliseconds: number = projectionDebounceMilliseconds
): void {
  if (
    room.document === null ||
    room.inFlightPersistenceBatch !== null ||
    room.pendingYjsUpdates.length > 0 ||
    room.pendingPersistenceUpdates.length > 0 ||
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
      room.inFlightPersistenceBatch !== null ||
      room.pendingYjsUpdates.length > 0 ||
      room.pendingPersistenceUpdates.length > 0 ||
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
      payload = Buffer.from(
        JSON.stringify({
          schemaId: "notezy.blocknote",
          schemaVersion: 1,
          projectedSequence,
          blocks: blockNoteProjector.projectYjsDocument(room.document),
        })
      );
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
    if (
      sendInternalFrame(
        subscriber.webSocket,
        InternalFrameType.InternalFrameType_ApplyBlockProjection,
        subscriber.connectionId,
        subscriber.connectorChannelId,
        blockPackId,
        payload
      )
    ) {
      return;
    }

    room.inFlightProjection = null;
    scheduleBlockProjection(room, blockPackId, projectionRetryMilliseconds);
  }, delayMilliseconds);
}

/* ============================== Batch handler ============================== */

const batchYjsHandler = new BatchYjsHandler(sendInternalFrame, resyncRoom);

/* ============================== HTTP server ============================== */

const server = createServer((request, response) => {
  const requestUrl = new URL(
    request.url ?? "/",
    `http://${request.headers.host ?? "localhost"}`
  );
  if (request.method === "GET" && requestUrl.pathname === "/healthz") {
    response.writeHead(200, { "content-type": "application/json" });
    response.end(
      JSON.stringify({ status: "ok", activeRoomCount: roomRegistry.size })
    );

    return;
  }

  response.writeHead(404);
  response.end();
});

/* ============================== WebSocket upgrade ============================== */

server.on("upgrade", (request, socket, head) => {
  const requestUrl = new URL(
    request.url ?? "/",
    `http://${request.headers.host ?? "localhost"}`
  );
  if (requestUrl.pathname !== "/internal/realtime/v1") {
    socket.destroy();

    return;
  }

  webSocketServer.handleUpgrade(request, socket, head, webSocket => {
    webSocketServer.emit("connection", webSocket, request);
  });
});

/* ============================== WebSocket connection ============================== */

webSocketServer.on("connection", webSocket => {
  webSocket.on("close", () => {
    roomRegistry.detachAll(webSocket);
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
          frame.connectorChannelId
        );
        if (room.document !== null) {
          if (room.inFlightPersistenceBatch !== null) {
            batchYjsHandler.retryInFlight(
              room,
              frame.blockPackId,
              webSocket,
              frame.connectionId,
              frame.connectorChannelId
            );

            return;
          }

          if (room.pendingPersistenceUpdates.length > 0) {
            batchYjsHandler.flush(
              room,
              frame.blockPackId,
              webSocket,
              frame.connectionId,
              frame.connectorChannelId
            );

            return;
          }

          sendRoomInitialState(room, frame.blockPackId);
          scheduleBlockProjection(room, frame.blockPackId);

          return;
        }

        if (room.isLoading) {
          return;
        }

        room.isLoading = true;
        if (
          sendInternalFrame(
            webSocket,
            InternalFrameType.InternalFrameType_LoadYjsDocument,
            frame.connectionId,
            frame.connectorChannelId,
            frame.blockPackId
          )
        ) {
          return;
        }

        resyncRoom(room, frame.blockPackId);

        return;
      }
      case InternalFrameType.InternalFrameType_Detach: {
        const room = roomRegistry.detach(
          frame.blockPackId,
          frame.connectionId,
          frame.connectorChannelId
        );
        if (room !== undefined && room.subscribers.size === 0) {
          batchYjsHandler.flush(room, frame.blockPackId);
        }

        return;
      }
      case InternalFrameType.InternalFrameType_YjsDocument: {
        const room = roomRegistry.getSubscriber(
          frame.blockPackId,
          frame.connectionId,
          frame.connectorChannelId
        );
        if (room === undefined) {
          sendInternalFrame(
            webSocket,
            InternalFrameType.InternalFrameType_ResyncRequired,
            frame.connectionId,
            frame.connectorChannelId,
            frame.blockPackId
          );

          return;
        }

        batchYjsHandler.queueUpdate(room, { webSocket, frame });

        return;
      }
      case InternalFrameType.InternalFrameType_Awareness: {
        const room = roomRegistry.getSubscriber(
          frame.blockPackId,
          frame.connectionId,
          frame.connectorChannelId
        );
        if (room === undefined) {
          sendInternalFrame(
            webSocket,
            InternalFrameType.InternalFrameType_ResyncRequired,
            frame.connectionId,
            frame.connectorChannelId,
            frame.blockPackId
          );

          return;
        }

        broadcastInternalFrame(
          room,
          frame.type,
          frame.blockPackId,
          frame.payload
        );

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
          batchYjsHandler.queueDeferredUpdates(room);
          if (
            room.inFlightPersistenceBatch === null &&
            room.pendingPersistenceUpdates.length === 0
          ) {
            sendRoomInitialState(room, frame.blockPackId);
          }
          scheduleBlockProjection(room, frame.blockPackId);
        } catch {
          resyncRoom(room, frame.blockPackId);
        }

        return;
      }
      case InternalFrameType.InternalFrameType_YjsUpdatePersisted: {
        const room = roomRegistry.get(frame.blockPackId);
        if (room === undefined) {
          return;
        }

        const inFlightPersistenceBatch = batchYjsHandler.handlePersisted(
          room,
          frame.blockPackId,
          frame.connectionId,
          frame.connectorChannelId,
          parseYjsUpdateSequence(frame.payload)
        );
        if (inFlightPersistenceBatch === null) {
          return;
        }

        broadcastInternalFrame(
          room,
          InternalFrameType.InternalFrameType_YjsDocument,
          frame.blockPackId,
          inFlightPersistenceBatch.payload
        );
        sendRoomInitialState(room, frame.blockPackId);
        batchYjsHandler.flush(room, frame.blockPackId);
        scheduleBlockProjection(room, frame.blockPackId);

        return;
      }
      case InternalFrameType.InternalFrameType_YjsPersistenceFailed: {
        const room = roomRegistry.get(frame.blockPackId);
        if (room !== undefined) {
          batchYjsHandler.handlePersistenceFailure(
            room,
            frame.blockPackId,
            frame.payload
          );
        }

        return;
      }
      case InternalFrameType.InternalFrameType_BlockProjectionApplied: {
        const room = roomRegistry.get(frame.blockPackId);
        let projectedUntilSequence: number | null = null;
        try {
          const value: unknown = JSON.parse(frame.payload.toString("utf8"));
          if (
            value !== null &&
            typeof value === "object" &&
            "projectedUntilSequence" in value &&
            typeof value.projectedUntilSequence === "number" &&
            Number.isSafeInteger(value.projectedUntilSequence) &&
            value.projectedUntilSequence >= -1
          ) {
            projectedUntilSequence = value.projectedUntilSequence;
          }
        } catch {}

        if (
          room === undefined ||
          projectedUntilSequence === null ||
          room.inFlightProjection === null ||
          room.inFlightProjection.connectionId !== frame.connectionId ||
          room.inFlightProjection.connectorChannelId !==
            frame.connectorChannelId ||
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
          room.inFlightProjection.connectorChannelId !==
            frame.connectorChannelId
        ) {
          return;
        }

        room.inFlightProjection = null;
        scheduleBlockProjection(
          room,
          frame.blockPackId,
          projectionRetryMilliseconds
        );

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

async function shutdown(signal: string): Promise<void> {
  console.info(`received ${signal}, stopping yjs worker`);

  server.close();
  for (const [blockPackId, room] of roomRegistry.entries()) {
    batchYjsHandler.flush(room, blockPackId);
  }

  const shutdownDeadline =
    Date.now() + YjsPersistenceBatchShutdownTimeoutMilliseconds;
  while (
    Date.now() < shutdownDeadline &&
    [...roomRegistry.entries()].some(
      ([, room]) =>
        room.inFlightPersistenceBatch !== null ||
        room.pendingPersistenceUpdates.length > 0
    )
  ) {
    await new Promise(resolve => setTimeout(resolve, 25));
  }

  webSocketServer.clients.forEach(webSocket => {
    webSocket.close(1001, "server shutdown");
  });
  webSocketServer.close();

  process.exit(0);
}

process.once("SIGINT", () => void shutdown("SIGINT"));
process.once("SIGTERM", () => void shutdown("SIGTERM"));
