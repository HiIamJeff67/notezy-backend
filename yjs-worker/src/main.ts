import { createServer } from "node:http";

import { WebSocket, WebSocketServer } from "ws";
import * as Y from "yjs";

import { config } from "./config.js";
import { createInternalFrame, parseInternalFrame } from "./internal_frame.js";
import { RoomRegistry } from "./room_registry.js";
import { internalRealtimePath, InternalFrameType } from "./types.js";

const roomRegistry = new RoomRegistry();
const webSocketServer = new WebSocketServer({ noServer: true });

/* ============================== Internal frame delivery ============================== */

function broadcastInternalFrame(
  room: ReturnType<RoomRegistry["getOrCreate"]>,
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
): void {
  if (webSocket.readyState !== WebSocket.OPEN) {
    return;
  }

  webSocket.send(
    createInternalFrame(type, connectionId, connectorChannelId, blockPackId, payload),
  );
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
  if (requestUrl.pathname !== internalRealtimePath) {
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
          frame.connectorChannelId,
        );
        sendInternalFrame(
          webSocket,
          InternalFrameType.InternalFrameType_YjsDocument,
          frame.connectionId,
          frame.connectorChannelId,
          frame.blockPackId,
          Buffer.from(Y.encodeStateAsUpdate(room.document)),
        );
        console.info("attached block pack room", {
          blockPackId: frame.blockPackId,
          connectionId: frame.connectionId,
          connectorChannelId: frame.connectorChannelId,
        });

        return;
      }
      case InternalFrameType.InternalFrameType_Detach:
        roomRegistry.detach(frame.blockPackId, frame.connectionId, frame.connectorChannelId);
        console.info("detached block pack room", {
          blockPackId: frame.blockPackId,
          connectionId: frame.connectionId,
          connectorChannelId: frame.connectorChannelId,
        });

        return;
      case InternalFrameType.InternalFrameType_YjsDocument:
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

        if (frame.type === InternalFrameType.InternalFrameType_YjsDocument) {
          try {
            Y.applyUpdate(room.document, frame.payload);
            room.dirtyUpdateCount += 1;
          } catch {
            sendInternalFrame(
              webSocket,
              InternalFrameType.InternalFrameType_ResyncRequired,
              frame.connectionId,
              frame.connectorChannelId,
              frame.blockPackId,
            );

            return;
          }
        }

        broadcastInternalFrame(room, frame.type, frame.blockPackId, frame.payload);

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
