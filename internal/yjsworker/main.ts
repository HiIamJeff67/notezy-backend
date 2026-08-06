import { SeverityNumber } from "@opentelemetry/api-logs";

import { YjsWorkerServer } from "./server.js";
import { Telemetry } from "./telemetry.js";

const telemetry = Telemetry.initialize();
const yjsWorkerServer = new YjsWorkerServer(telemetry);
let isShuttingDown = false;

async function shutdown(signal: string): Promise<void> {
  if (isShuttingDown) {
    return;
  }
  isShuttingDown = true;

  telemetry.log(SeverityNumber.INFO, "yjs_worker.shutdown_started", { signal });
  await yjsWorkerServer.shutdown();
  await telemetry.shutdown();

  process.exit(0);
}

process.once("SIGINT", () => void shutdown("SIGINT"));
process.once("SIGTERM", () => void shutdown("SIGTERM"));
