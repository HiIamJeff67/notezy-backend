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

export function convertBytesToUUIDString(bytes: Buffer): string | null {
  if (bytes.length !== 16 || bytes.every((byte) => byte === 0)) {
    return null;
  }
  const hex = bytes.toString("hex");
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}

export function convertUUIDToBytes(uuidString: string): Buffer {
  if (!/^[0-9a-f]{8}-[0-9a-f]{4}-[1-8][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(uuidString)) {
    throw new Error("invalid UUID");
  }
  return Buffer.from(uuidString.replaceAll("-", ""), "hex")
}
