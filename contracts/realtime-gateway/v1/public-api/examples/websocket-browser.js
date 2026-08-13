// Paste in a browser console after obtaining tickets from Gateway v1.
// Set these values locally; never commit a real ticket.
const realtimeEndpoint = "ws://localhost/realtime/development/v1";
const connectionTicket = "<single-use-connection-ticket>";
const channelTicket = "<single-use-channel-ticket>";
const blockPackId = "00000000-0000-4000-8000-000000000001";

const socket = new WebSocket(realtimeEndpoint, connectionTicket);
socket.binaryType = "arraybuffer";
socket.addEventListener("message", ({ data }) => {
  if (typeof data === "string") {
    const frame = JSON.parse(data);
    console.log("control", frame);
    if (frame.type === "ready") {
      socket.send(JSON.stringify({ version: 1, type: "subscribe", requestId: "sub-1", channelType: "BlockPack", channelId: blockPackId, channelTicket }));
    }
    return;
  }
  const bytes = new Uint8Array(data);
  const connectorChannelId = new DataView(data).getUint32(2, false);
  console.log("binary", { version: bytes[0], type: bytes[1], connectorChannelId, payload: bytes.slice(6) });
});
socket.addEventListener("close", ({ code, reason }) => console.log("closed", code, reason));
socket.addEventListener("error", (error) => console.error("websocket error", error));
