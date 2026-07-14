import { upgradeWebSocket } from "@hono/node-server";
import type { Hono } from "hono";
import type WebSocket from "ws";

export function configureRealtimeRoutes(
  app: Hono,
  handleConnection: (webSocket: WebSocket) => void
): void {
  app.get(
    "/internal/realtime/v1",
    upgradeWebSocket(() => ({
      onOpen(_event, webSocket) {
        handleConnection(webSocket.raw as WebSocket);
      },
    }))
  );
}
