import { createServer } from "node:http";

import { type RawData, WebSocketServer } from "ws";

import { config } from "./config.js";
import { internalRealtimePath, InternalFrameType } from "./types.js";
import { parseInternalFrame } from "./internal_frame.js";
import { RoomRegistry } from "./room_registry.js";

const roomRegistry = new RoomRegistry();
const webSocketServer = new WebSocketServer({ noServer: true });
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

webSocketServer.on("connection", (webSocket) => {
  webSocket.on("message", (payload, isBinary) => {
    if (!isBinary) {
      webSocket.close(1003, "internal realtime frames must be binary");

      return;
    }

    const frame = parseInternalFrame(rawDataToBuffer(payload));
    if (frame === null || frame.version !== 1) {
      webSocket.close(1002, "invalid internal realtime frame");

      return;
    }

    if (frame.type === InternalFrameType.InternalFrameType_Attach) {
      roomRegistry.getOrCreate(frame.blockPackId);
      console.info("attached block pack room", {
        blockPackId: frame.blockPackId,
        connectionId: frame.connectionId,
        connectorChannelId: frame.connectorChannelId,
      });

      return;
    }

    if (frame.type === InternalFrameType.InternalFrameType_Detach) {
      console.info("detached block pack room", {
        blockPackId: frame.blockPackId,
        connectionId: frame.connectionId,
        connectorChannelId: frame.connectorChannelId,
      });

      return;
    }

    console.warn("received internal frame before its handler is enabled", {
      type: frame.type,
      blockPackId: frame.blockPackId,
    });
  });
});

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

function rawDataToBuffer(payload: RawData): Buffer {
  if (Buffer.isBuffer(payload)) {
    return payload;
  }
  if (payload instanceof ArrayBuffer) {
    return Buffer.from(payload);
  }
  if (Array.isArray(payload)) {
    return Buffer.concat(payload);
  }

  throw new Error("unsupported websocket raw data");
}
