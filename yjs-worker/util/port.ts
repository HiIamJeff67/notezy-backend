export function parsePort(portString: string | undefined): number {
  if (portString === undefined || portString === "") {
    return 8787;
  }

  const port = Number(portString);
  if (!Number.isInteger(port) || port < 1 || port > 65535) {
    throw new Error("YJS_WORKER_PORT must be an integer between 1 and 65535");
  }

  return port;
}
