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

The BlockPack response contains `channelTicket`, `expiresAt`, `channelType`, `channelId`, `permission`, `roomName`, `fragmentName`, `schemaId`, `schemaVersion`, `realtimeProtocolVersion`, `lastUpdateSequence`, and `compactedUntilSequence`. A BlockPackYjsDocument row is created in the same transaction as its BlockPack; a missing row is a backend error, not an empty document state.

`permission: "read"` is available to Read, Write, Admin, and Owner users. `permission: "write"` is available only to Write, Admin, and Owner users. A BlockPack, its Yjs document, parent SubShelf, and RootShelf must all be active before a ticket is issued.

## Realtime Participants

Owners and Admins can inspect the currently active connections in one BlockPack room:

```text
GET /api/development/v1/realtime/blockPacks/:blockPackId/participants
```

The response body uses the normal `{ success, data, exception }` envelope. `data` is an array of:

```json
{
  "userPublicId": "UUID",
  "name": "string",
  "displayName": "string",
  "channelPermission": "read" | "write",
  "connectionCount": 1
}
```

Participants are derived from active Redis subscriber leases. They are ephemeral presence data, not an access-control source; an empty array means no active room connection was observed.

Tickets are EdDSA JWTs signed and verified by Go. Go receives `REALTIME_TICKET_PRIVATE_KEY_BASE64`, which is Base64-encoded PKCS#8 Ed25519 DER. Node worker 不接收 ticket key，也不驗證 public ticket；它只接受已由 Go Gateway 驗證後送出的 internal frame。Tickets contain `iss`, `aud`, `sub`, `jti`, `iat`, `exp`, a hash of the `User-Agent`, and the channel claims where applicable. Audiences are `notezy-realtime-connection` and `notezy-realtime-block-pack`.

Generate the two deployment values once and store them in secret management, never in the repository:

```bash
openssl genpkey -algorithm ED25519 -out realtime-ticket-private.pem
REALTIME_TICKET_PRIVATE_KEY_BASE64="$(openssl pkcs8 -topk8 -nocrypt -in realtime-ticket-private.pem -outform DER | base64 | tr -d '\n')"
```

Tickets are short-lived for five minutes and stateless. `jti` is a trace identifier; it is not a one-time-use guarantee. True replay prevention would require shared state and is intentionally not introduced in this phase.

`NOT-7` replaces root-upgrade access-token middleware validation with the connection ticket sent as the single `Sec-WebSocket-Protocol` value. Client 建立 socket 時傳入 `new WebSocket(realtimeEndpoint, [connectionTicket])`；server 驗證 signed `User-Agent` hash 後選擇同一 subprotocol。每一個 subscribe 都必須再帶入並驗證自己的 `channelTicket`。connection 與 channel ticket 的 `sub` 必須一致。

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

Gateway 僅將 subscribed Yjs/awareness frame 轉送至已分派的 worker。`yjs-document` 是 mutation，channel ticket 必須具有 `permission: "write"`；read channel 收到 `channel_permission_denied`，且 payload 不會送入 worker。awareness 可由 read/write channel 使用，但由 worker 維護 room-level `y-protocols/awareness` state；它不寫入 `Y.Doc`、snapshot 或 durable update log。payload 是 `encodeAwarenessUpdate` 的 raw encoded update，不包裝 y-websocket protocol header。worker 驗證每個 client ID 僅能被一個 connector channel 宣告，拒絕 malformed、重複或跨 connector 的 client ID，並從 authoritative state 重新編碼後 relay。attach 成功且 document initial state 已送出後，worker 會補送目前 awareness state；detach、socket close、worker resync 會 broadcast 對應 client ID 的 null-state removal，避免 ghost cursor/presence。`yjs-document` payload 是 `Y.encodeStateAsUpdate` 產生的 raw encoded update。後續每個有效 update 都先套用到 memory `Y.Doc`、append 至 durable update log，收到 append ACK 後才轉送給同一 BlockPack 的所有 subscriber。Yjs update 可重複套用，因此 sender 收到自己的 relay 不影響正確性。

## Errors And Future Lifecycle

All gateway errors are JSON:

```json
{ "version": 1, "type": "error", "requestId": "sub-1", "connectorChannelId": 1, "code": "channel_not_found", "message": "connectorChannelId is not subscribed on this connection" }
```

Stable error codes are `authentication_managed_by_upgrade`, `binary_channel_not_ready`, `channel_backpressure`, `channel_limit_exceeded`, `channel_not_found`, `channel_permission_denied`, `invalid_acknowledgement`, `invalid_binary_frame`, `invalid_channel_id`, `invalid_channel_ticket`, `invalid_channel_type`, `invalid_connector_channel_id`, `invalid_control_frame`, `permission_revoked`, `resource_unavailable`, `room_admission_unavailable`, `room_connection_limit_exceeded`, `resubscribe_required`, `unsupported_binary_type`, `unsupported_channel_type`, `unsupported_control_type`, `unsupported_message_type`, `unsupported_protocol_version`, and `worker_unavailable`.

`permission_revoked` means the user no longer has the requested channel permission. `resource_unavailable` means the BlockPack, Yjs document, SubShelf, or RootShelf is deleted or otherwise inactive. Both errors remove the logical channel; the client must stop its editor/provider and must not continue sending binary frames on that `connectorChannelId`.

Before upgrading a public socket, Go applies an IP-based upgrade rate limit, its per-process connector cap, and a distributed per-user root-connection cap. A rejected upgrade returns HTTP `429` for the user cap or HTTP `503` for gateway/admission availability; the client must not retry in a tight loop. The distributed cap is represented by Redis TTL leases and is refreshed with the transport heartbeat, so an abnormal process exit is recovered automatically after the lease expires.

Each BlockPack subscription is re-authorized at subscribe time and consumes one shared active-subscriber slot from the BlockPack owner's plan. Before a `yjs-document` mutation is forwarded, Go revalidates the active BlockPack hierarchy and the writer permission. A validation failure releases the lease, detaches the worker channel, and returns `permission_revoked` or `resource_unavailable`. Read and write subscriptions use the same room capacity. `room_connection_limit_exceeded` means the client should close or unsubscribe another active subscriber before retrying. A successful `unsubscribe`, connection close, permission revocation, and lease expiry all release the slot.

This monolith phase does not yet fan out an external permission/delete event to every Gateway process. An idle channel can therefore remain open until it next subscribes or submits a document mutation. Kafka-backed lifecycle fanout is intentionally deferred to the microservice architecture project.

The gateway caps a connection at 64 active channels. Released IDs are not reused during that connection. Public outbound data uses a bounded queue per `connectorChannelId`, with round-robin scheduling between channels. JSON control frames are always scheduled before binary data. Each channel allows at most 256 queued binary frames and 4 MiB of queued binary payload. Awareness is ephemeral: a queued awareness frame replaces the previous queued awareness frame for that channel. Yjs document updates are never silently dropped or coalesced by Go; if their channel queue is full, the gateway detaches only that channel and sends `channel_backpressure`, requiring a resubscribe/resync while unrelated channels remain active. A failed read or a write that cannot complete within 10 seconds closes the physical socket. Go-to-worker multiplexing uses `YJS_WORKER_URLS`, a comma-separated internal endpoint list. Each `blockPackId` maps consistently to one endpoint; each endpoint has one long-lived internal WebSocket and a bounded outbound queue. An unavailable worker or a full internal queue rejects the affected channel payload with `worker_unavailable`.

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

Internal types are `1` `attach`, `2` `detach`, `3` `yjs-document`, `4` `awareness`, `5` `resync-required`, `6` `permission-revoked`, `7` `load-yjs-document`, `8` `yjs-document-loaded`, `9` `append-yjs-update` (legacy single update), `10` `yjs-update-persisted`, `11` `yjs-persistence-failed`, `12` `apply-block-projection`, `13` `block-projection-applied`, `14` `block-projection-failed`, and `15` `append-yjs-update-batch`.

`attach` and `detach` are idempotent. A first attach asks Go for a binary cold-load payload: `lastUpdateSequence(int64)`, `compactedUntilSequence(int64)`, `projectedUntilSequence(int64)`, snapshot length/state-vector length/update count (`uint32` each), snapshot bytes, state-vector bytes, then ordered update entries of `updateSequence(int64)`, payload length (`uint32`), raw update bytes. The worker materializes the document before it sends the public initial state.

`append-yjs-update-batch` carries `[persistenceBatchId:16][originConnectionId:16, zero UUID when mixed][merged raw Yjs update:n]`. Go appends it transactionally and returns `yjs-update-persisted` with its `updateSequence(int64)` payload. The worker serializes persistence batches per BlockPack and broadcasts only after this ACK. `persistenceBatchId` makes a retry idempotent when PostgreSQL committed but the ACK was lost. On an internal worker reconnect, Go replays `attach` for every active channel assigned to that worker before it forwards a client payload. When replay cannot be completed, Go emits `resync_required` to that public channel and waits for the client to resubscribe; it never silently broadcasts an unpersisted Yjs payload.

`apply-block-projection` carries UTF-8 JSON `{ schemaId, schemaVersion, projectedSequence, blocks }`; the BlockPack id is the internal frame `channelId`. This request is accepted only over Go-established private worker connections, not through public WebSocket or REST routes. Go validates the schema and durable sequence, bulk applies the BlockTable projection, and returns JSON `{ applied, projectedUntilSequence }` with `block-projection-applied`; malformed, stale-invalid, or failed requests receive `block-projection-failed`.

The internal implementation uses a bounded outbound queue per worker. Queue exhaustion or a dead worker closes affected logical channels with `worker_unavailable`; public sockets may remain open for their unrelated channels. The public per-channel queue and backpressure behavior belongs to `NOT-22`; the envelope and recovery semantics above remain the protocol contract.
