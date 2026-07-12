# Realtime Protocol Contract

Phase 0 endpoint:

| Environment | URL |
| --- | --- |
| local through nginx | `ws://localhost/realtime/development/v1` |
| production | `wss://api.notezy.app/realtime/development/v1` |

A physical WebSocket belongs to one client app instance. Each new connection receives a new `connectionId`, and its `connectorChannelId` values are valid only for that connection.

## Realtime Ticket APIs

`NOT-5` adds two authenticated REST endpoints below the API base path. They use the normal access-token/cookie authentication pipeline and never accept an access token in the WebSocket URL.

| Method | URL | Body | Purpose |
| --- | --- | --- | --- |
| `POST` | `/api/development/v1/realtime/createMyRealtimeConnectionTicket` | none | Issue a root connection ticket. |
| `POST` | `/api/development/v1/realtime/createMyBlockPackChannelTicket` | `{ "blockPackId": "UUID", "permission": "read" \| "write" }` | Check the current user's non-deleted BlockPack permission and issue a capability for that one BlockPack. |

The connection response contains `realtimeEndpoint` (`/realtime/development/v1`), `realtimeProtocolVersion`, `connectionTicket`, and `expiresAt`.

The BlockPack response contains `channelTicket`, `expiresAt`, `channelType`, `channelId`, `permission`, `roomName`, `fragmentName`, `schemaId`, `schemaVersion`, `realtimeProtocolVersion`, `lastUpdateSequence`, and `compactedUntilSequence`. When the Yjs document row is not created yet, both returned sequences are `0`.

`permission: "read"` is available to Read, Write, Admin, and Owner users. `permission: "write"` is available only to Write, Admin, and Owner users. Soft-deleted BlockPacks do not receive a ticket.

Tickets are EdDSA JWTs signed by Go. Go receives `REALTIME_TICKET_PRIVATE_KEY_BASE64`, which is Base64-encoded PKCS#8 Ed25519 DER. The future Node worker receives the matching `REALTIME_TICKET_PUBLIC_KEY_BASE64`, which is Base64-encoded SPKI Ed25519 DER. Tickets contain `iss`, `aud`, `sub`, `jti`, `iat`, `exp`, a hash of the `User-Agent`, and the channel claims where applicable. Audiences are `notezy-realtime-connection` and `notezy-realtime-block-pack`.

Generate the two deployment values once and store them in secret management, never in the repository:

```bash
openssl genpkey -algorithm ED25519 -out realtime-ticket-private.pem
REALTIME_TICKET_PRIVATE_KEY_BASE64="$(openssl pkcs8 -topk8 -nocrypt -in realtime-ticket-private.pem -outform DER | base64 | tr -d '\n')"
REALTIME_TICKET_PUBLIC_KEY_BASE64="$(openssl pkey -in realtime-ticket-private.pem -pubout -outform DER | base64 | tr -d '\n')"
```

Tickets are short-lived for five minutes and stateless. `jti` is a trace identifier; it is not a one-time-use guarantee. True replay prevention would require shared state and is intentionally not introduced in this phase.

The existing Phase 0 gateway still authenticates the upgrade through its current middleware and does not yet validate tickets. `NOT-7` replaces that root-upgrade validation with the connection ticket sent through `Sec-WebSocket-Protocol`, validates the signed `User-Agent` hash against the upgrade request, then validates `channelTicket` on every subscribe request. This temporary state is only for completing the backend boundary without breaking the current frontend smoke client.

## Text Control Frames

All control frames are UTF-8 JSON and begin with `version`, `type`, and an optional client-generated `requestId`. The current version is `1`.

```json
{ "version": 1, "type": "subscribe", "requestId": "sub-1", "channelType": "BlockPack", "channelId": "4b49c1fc-8c68-40da-84b5-c5808201504a", "channelTicket": "<channel ticket>" }
```

`channelType` and `channelId` identify the resource. `connectorChannelId` is the unsigned connection-local ID used in binary frames, `ack`, and `unsubscribe`. Repeating the same `channelType + channelId` subscription is idempotent and returns the same `connectorChannelId` with `existing: true`.

Phase 0 enables only `channelType: "BlockPack"`; other values receive `unsupported_channel_type`. Adding a new channel type requires one explicit `subscribe` branch, an internal channel-type code, and its own capability/worker handling.

```json
{ "version": 1, "type": "subscribed", "requestId": "sub-1", "channelType": "BlockPack", "channelId": "4b49c1fc-8c68-40da-84b5-c5808201504a", "connectorChannelId": 1, "existing": false }
```

```json
{ "version": 1, "type": "unsubscribe", "requestId": "unsub-1", "connectorChannelId": 1 }
```

```json
{ "version": 1, "type": "unsubscribed", "requestId": "unsub-1", "channelType": "BlockPack", "channelId": "4b49c1fc-8c68-40da-84b5-c5808201504a", "connectorChannelId": 1 }
```

`ack` advances the client-confirmed sequence for a channel. Its sequence must never move backwards.

```json
{ "version": 1, "type": "ack", "requestId": "ack-1", "connectorChannelId": 1, "sequence": 42 }
{ "version": 1, "type": "acknowledged", "requestId": "ack-1", "connectorChannelId": 1, "sequence": 42 }
```

`ping` returns `pong`. `heartbeat` returns a `heartbeat` with `unixMilliNow`; native WebSocket ping/pong is also used by the gateway to keep the transport alive. A client must treat a new `ready` frame as a reconnect boundary and subscribe every required BlockPack again.

```json
{ "version": 1, "type": "ready", "connectionId": "d3eaa2e9-bb1a-4b6b-af5d-e4f102b27b62", "resubscribeRequired": true }
```

`authenticate` is deliberately rejected with `authentication_managed_by_upgrade`; root connection authentication is not a channel operation. `channelTicket` is present in the subscribe envelope now and is enforced by `NOT-7`.

## Binary Frames

Binary frames never Base64-encode Yjs data and never use JSON block events. Their header is exactly six bytes, followed by raw bytes:

| Offset | Length | Value |
| --- | --- | --- |
| `0` | 1 byte | protocol version (`1`) |
| `1` | 1 byte | binary type: `1` = `yjs-document`, `2` = `awareness` |
| `2` | 4 bytes | unsigned big-endian `connectorChannelId` |
| `6` | remaining | raw Yjs or awareness payload |

The `connectorChannelId` maps the payload to its subscribed `channelType + channelId`; a public binary frame therefore does not repeat the resource identity. Unknown, unsubscribed, malformed, or unsupported binary frames receive an error JSON frame and are never forwarded.

Phase 0 validates the header and channel lifecycle only. It responds with `binary_channel_not_ready` after validating a subscribed Yjs/awareness frame because worker forwarding begins in `NOT-7`; no payload is persisted or broadcast yet.

## Errors And Future Lifecycle

All gateway errors are JSON:

```json
{ "version": 1, "type": "error", "requestId": "sub-1", "connectorChannelId": 1, "code": "channel_not_found", "message": "connectorChannelId is not subscribed on this connection" }
```

Stable Phase 0 error codes are `authentication_managed_by_upgrade`, `binary_channel_not_ready`, `channel_limit_exceeded`, `channel_not_found`, `invalid_acknowledgement`, `invalid_binary_frame`, `invalid_channel_id`, `invalid_channel_type`, `invalid_connector_channel_id`, `invalid_control_frame`, `permission_revoked`, `resubscribe_required`, `unsupported_binary_type`, `unsupported_channel_type`, `unsupported_control_type`, `unsupported_message_type`, `unsupported_protocol_version`, and `worker_unavailable`.

The gateway caps a connection at 64 active channels. Released IDs are not reused during that connection. Public writes are serialized without an unbounded queue; a failed read or a write that cannot complete within 10 seconds closes the physical socket. `permission_revoked`, capability-ticket validation, Go-to-worker multiplexing, and worker forwarding are enabled by their follow-up issues; they must preserve this wire header and channel-ID ownership model.

## Internal Go Gateway To Yjs Worker Frames

The future Go-to-worker transport is a small pool of long-lived multiplex WebSockets per Node worker. It must never create one internal WebSocket per public client. Its binary frame header is fixed now so Go and Node can implement it independently:

| Offset | Length | Value |
| --- | --- | --- |
| `0` | 1 byte | worker protocol version (`1`) |
| `1` | 1 byte | internal type |
| `2` | 1 byte | internal channel-type code: `1` = `BlockPack` |
| `3` | 16 bytes | raw UUID `connectionId` |
| `19` | 4 bytes | unsigned big-endian `connectorChannelId` |
| `23` | 16 bytes | raw UUID `channelId` |
| `39` | remaining | raw Yjs/awareness payload, or empty for attach/detach |

Internal types are `1` `attach`, `2` `detach`, `3` `yjs-document`, `4` `awareness`, `5` `resync-required`, and `6` `permission-revoked`. `attach` and `detach` are idempotent. On an internal worker reconnect, Go replays `attach` for every active channel assigned to that worker before it forwards a client payload. When replay cannot be completed, Go emits `resync_required` to that public channel and waits for the client to resubscribe; it never silently drops an accepted Yjs payload.

The internal implementation uses a bounded outbound queue per worker. Queue exhaustion or a dead worker closes affected logical channels with `worker_unavailable`; public sockets may remain open for their unrelated channels. Those execution details belong to `NOT-7`, while the envelope and recovery semantics above are the Phase 0 contract.
