import type { Hono } from "hono";

export function configureHealthRoutes(
  app: Hono,
  getActiveRoomCount: () => number
): void {
  app.get("/healthz", context =>
    context.json({ status: "ok", activeRoomCount: getActiveRoomCount() })
  );
}
