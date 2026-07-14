import { YjsWorkerServer } from "./server.js";

const yjsWorkerServer = new YjsWorkerServer();
let isShuttingDown = false;

async function shutdown(signal: string): Promise<void> {
  if (isShuttingDown) {
    return;
  }
  isShuttingDown = true;

  console.info(`received ${signal}, stopping yjs worker`);
  await yjsWorkerServer.shutdown();

  process.exit(0);
}

process.once("SIGINT", () => void shutdown("SIGINT"));
process.once("SIGTERM", () => void shutdown("SIGTERM"));
