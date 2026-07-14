import { serve } from "@hono/node-server";
import { Hono } from "hono";
import { WebSocketServer } from "ws";

import { config } from "./config.js";
import { BlockPackProjector } from "./realtime/block_pack_projector.js";
import { RealtimeGateway } from "./realtime/gateway.js";
import { RoomRegistry } from "./realtime/room_registry.js";
import { YjsDebouncer } from "./realtime/yjs_debouncer.js";
import { configureHealthRoutes } from "./routes/health_route.js";
import { configureRealtimeRoutes } from "./routes/realtime_route.js";
import { configureYjsCompactionRoutes } from "./routes/yjs_compaction_route.js";
import { YjsCompactionService } from "./services/yjs_compaction_service.js";

export class YjsWorkerServer {
  private readonly server: ReturnType<typeof serve>;
  private readonly webSocketServer: WebSocketServer;
  private readonly realtimeGateway: RealtimeGateway;

  constructor() {
    const app = new Hono();
    const yjsCompactionService = new YjsCompactionService();
    this.webSocketServer = new WebSocketServer({ noServer: true });
    this.realtimeGateway = new RealtimeGateway(
      new RoomRegistry(),
      new BlockPackProjector(),
      yjsCompactionService,
      new YjsDebouncer()
    );

    configureHealthRoutes(
      app,
      this.realtimeGateway.getActiveRoomCount.bind(this.realtimeGateway)
    );
    configureYjsCompactionRoutes(app, yjsCompactionService);
    configureRealtimeRoutes(
      app,
      this.realtimeGateway.handleConnection.bind(this.realtimeGateway)
    );

    this.server = serve(
      {
        fetch: app.fetch,
        hostname: config.host,
        port: config.port,
        websocket: { server: this.webSocketServer },
      },
      () => {
        console.info(`yjs worker listening on ${config.host}:${config.port}`);
      }
    );
  }

  async shutdown(): Promise<void> {
    const closeServer = new Promise<void>(resolve => {
      this.server.close(() => resolve());
    });

    await this.realtimeGateway.shutdown();
    await closeServer;
  }
}
