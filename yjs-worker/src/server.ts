import { serve } from "@hono/node-server";
import { SeverityNumber } from "@opentelemetry/api-logs";
import { Hono } from "hono";
import { WebSocketServer } from "ws";

import { config } from "./config.js";
import { BlockPackProjector } from "./realtime/block_pack_projector.js";
import { RealtimeGateway } from "./realtime/gateway.js";
import { RoomRegistry } from "./realtime/room_registry.js";
import { configureHealthRoutes } from "./routes/health_route.js";
import { configureRealtimeRoutes } from "./routes/realtime_route.js";
import { configureYjsCompactionRoutes } from "./routes/yjs_compaction_route.js";
import { configureYjsProjectionRoutes } from "./routes/yjs_projection_route.js";
import { YjsCompactionService } from "./services/yjs_compaction_service.js";
import { YjsProjectionService } from "./services/yjs_projection_service.js";
import type { Telemetry } from "./telemetry.js";

export class YjsWorkerServer {
  private readonly server: ReturnType<typeof serve>;
  private readonly webSocketServer: WebSocketServer;
  private readonly realtimeGateway: RealtimeGateway;

  constructor(telemetry: Telemetry) {
    const app = new Hono();
    const blockPackProjector = new BlockPackProjector();
    const yjsCompactionService = new YjsCompactionService(telemetry);
    const yjsProjectionService = new YjsProjectionService(
      blockPackProjector,
      telemetry
    );
    const roomRegistry = new RoomRegistry(telemetry);
    this.webSocketServer = new WebSocketServer({ noServer: true });
    this.realtimeGateway = new RealtimeGateway(
      roomRegistry,
      yjsCompactionService,
      telemetry
    );

    configureHealthRoutes(
      app,
      this.realtimeGateway.getActiveRoomCount.bind(this.realtimeGateway)
    );
    configureRealtimeRoutes(
      app,
      this.realtimeGateway.handleConnection.bind(this.realtimeGateway)
    );
    configureYjsCompactionRoutes(app, yjsCompactionService, telemetry);
    configureYjsProjectionRoutes(app, yjsProjectionService, telemetry);

    this.server = serve(
      {
        fetch: app.fetch,
        hostname: config.host,
        port: config.port,
        websocket: { server: this.webSocketServer },
      },
      () => {
        telemetry.log(SeverityNumber.INFO, "yjs_worker.started", {
          host: config.host,
          port: config.port,
        });
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
