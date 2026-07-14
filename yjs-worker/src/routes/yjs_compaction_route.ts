import type { Hono } from "hono";
import { bodyLimit } from "hono/body-limit";

import { YjsMaintenanceMaximumPayloadBytes } from "../constants/yjs_maintenance.js";
import type { YjsCompactionService } from "../services/yjs_compaction_service.js";

export function configureYjsCompactionRoutes(
  app: Hono,
  yjsCompactionService: YjsCompactionService
): void {
  app.post(
    "/internal/yjs-compaction/v1",
    bodyLimit({
      maxSize: YjsMaintenanceMaximumPayloadBytes,
      onError: context => context.body(null, 413),
    }),
    async context => {
      const contentLength = Number(context.req.header("content-length") ?? 0);
      if (
        !Number.isSafeInteger(contentLength) ||
        contentLength <= 0 ||
        contentLength > YjsMaintenanceMaximumPayloadBytes
      ) {
        return context.body(null, 413);
      }

      const payload = Buffer.from(await context.req.arrayBuffer());

      try {
        return context.body(
          Uint8Array.from(yjsCompactionService.compactBatch(payload)),
          200,
          {
            "content-type": "application/octet-stream",
          }
        );
      } catch {
        return context.body(null, 422);
      }
    }
  );
}
