# Realtime Protocol Design

## Scope

This document is the backend-owned specification for the public Realtime
WebSocket protocol and the HTTP ticket APIs that authorize it. It defines the
wire contract shared by the Gateway, Yjs worker, persistence path, and any
protocol consumer. Resource-specific ticket behavior is documented in
[Realtime Editor API Design](../api-route-design/realtime-editor-api.md).

Phase 0 endpoint:

| Environment | URL |
| --- | --- |
| local through nginx | `ws://localhost/realtime/development/v1` |
| production | `wss://api.notezy.app/realtime/development/v1` |

A physical WebSocket belongs to one client app instance. Each new connection receives a new `connectionId`, and its `connectorChannelId` values are valid only for that connection.

## Realtime Ticket APIs

The authenticated REST endpoints below the API base path issue capabilities for
the public socket. They use the normal access-token/cookie authentication
pipeline and never accept an access token in the WebSocket URL.

| Method | URL | Body | Purpose |
| --- | --- | --- | --- |
| `POST` | `/api/development/v1/realtime/connection/ticket` | none | Issue a root connection ticket. |
| `POST` | `/api/development/v1/realtime/channel/block-pack/ticket` | `{ "blockPackId": "UUID", "permission": "read" \| "write" }` | Check the current user's non-deleted BlockPack permission and issue a capability for that one BlockPack. |

The connection response contains `realtimeEndpoint` (`/realtime/development/v1`), `realtimeProtocolVersion`, `connectionTicket`, and `expiresAt`.

The BlockPack response contains `channelTicket`, `expiresAt`, `channelType`, `channelId`, `permission`, `roomName`, `fragmentName`, `schemaId`, `schemaVersion`, `realtimeProtocolVersion`, `documentQuotaPolicyVersion`, `maximumBlockCount`, `lastUpdateSequence`, and `compactedUntilSequence`. The quota fields let the client prevent obviously invalid edits, but they are not an authorization boundary. A BlockPackYjsDocument row is created in the same transaction as its BlockPack; a missing row is a backend error, not an empty document state.

`permission: "read"` is available to Read, Write, Admin, and Owner users. `permission: "write"` is available only to Write, Admin, and Owner users. A BlockPack, its Yjs document, parent SubShelf, and RootShelf must all be active before a ticket is issued.

## Realtime Participants

Read, Write, Admin, and Owner users can inspect the currently active connections in one BlockPack room:

```text
GET /realtime/development/v1/block-pack/:blockPackId/participants
```

The response body uses the RealtimeGateway response envelope. `data.participants` is an array of ephemeral lease records:

```json
{
  "userPublicId": "UUID",
  "channelPermission": "read" | "write",
  "connectionCount": 1
}
```

RealtimeGateway derives participants directly from its active Redis subscriber
leases and serves the snapshot from its own public realtime endpoint. API Gateway
does not proxy this request and Core does not read RealtimeGateway Redis participant
state. Participants are ephemeral presence data, not an access-control source; an
empty array means no active room connection was observed. User profile details, if
needed by a client, remain a separate Core API query.

Tickets are EdDSA JWTs signed by Core and verified by RealtimeGateway. Core exclusively receives `REALTIME_TICKET_PRIVATE_KEY_BASE64`, which is Base64-encoded PKCS#8 Ed25519 DER. RealtimeGateway exclusively receives `REALTIME_TICKET_PUBLIC_KEY_BASE64`, which is Base64-encoded PKIX Ed25519 DER. The Node worker receives neither key and does not validate public tickets; it only accepts internal frames whose claims were verified by RealtimeGateway. Tickets contain `iss`, `aud`, `sub`, `jti`, `iat`, `exp`, a hash of the `User-Agent`, and the channel claims where applicable. Audiences are `notezy-realtime-connection` and `notezy-realtime-block-pack`.

Generate the two deployment values once and store them in secret management, never in the repository:

```bash
openssl genpkey -algorithm ED25519 -out realtime-ticket-private.pem
REALTIME_TICKET_PRIVATE_KEY_BASE64="$(openssl pkcs8 -topk8 -nocrypt -in realtime-ticket-private.pem -outform DER | base64 | tr -d '\n')"
REALTIME_TICKET_PUBLIC_KEY_BASE64="$(openssl pkey -in realtime-ticket-private.pem -pubout -outform DER | base64 | tr -d '\n')"
```

The private key is injected only into Core. The public key is injected only into
RealtimeGateway. A client may decode ticket claims, but modifying any protected
header or claim invalidates the Ed25519 signature.

Tickets are short-lived for five minutes and single-use. After signature and
claim validation, RealtimeGateway atomically consumes the `jti` with Redis
`SET NX` and a TTL ending at `exp`; a second use is rejected. This shared state
works across all RealtimeGateway instances. A failed connection or subscription
must request a newly issued ticket from Core.

Root-upgrade authentication uses the connection ticket as the sole
`Sec-WebSocket-Protocol` value. RealtimeGateway validates the signed `User-Agent`
hash and selects the same subprotocol. Every `subscribe` carries and validates
its own `channelTicket`; connection and channel ticket `sub` claims must
match.

## Text Control Frames

All control frames are UTF-8 JSON and begin with `version`, `type`, and an optional client-generated `requestId`. The current version is `1`.

```json
{ "version": 1, "type": "subscribe", "requestId": "sub-1", "channelType": "BlockPack", "channelId": "4b49c1fc-8c68-40da-84b5-c5808201504a", "channelTicket": "<channel ticket>" }
```

`channelType` and `channelId` identify the resource. `connectorChannelId` is the unsigned connection-local ID used in binary frames, `ack`, and `unsubscribe`. Repeating the same `channelType + channelId` subscription is idempotent and returns the same `connectorChannelId` with `existing: true`.

Phase 0 enables only `channelType: "BlockPack"`; other values receive `unsupported_channel_type`. Adding a new channel type requires one explicit `subscribe` branch, an internal channel-type code, and its own capability/worker handling.

```json
{ "version": 1, "type": "subscribed", "requestId": "sub-1", "channelType": "BlockPack", "channelId": "4b49c1fc-8c68-40da-84b5-c5808201504a", "connectorChannelId": 1, "existing": false, "documentQuotaPolicyVersion": 1, "maximumBlockCount": 1000, "participants": [{ "userPublicId": "UUID", "channelPermission": "read", "connectionCount": 1 }] }
```

`subscribed.participants` is the room presence snapshot at admission time. It
contains only public user IDs, current channel permission, and active connection
count; it never contains internal IDs, cookies, delegation data, or user profile
fields. `subscribed` is emitted only after Yjs Worker confirms that the subscriber
is registered and the authoritative document is ready; the client must not send
binary document or awareness frames before it. A repeated subscription returns
the same snapshot with `existing: true`.

After subscription, a client may receive these unsolicited room-scoped deltas:
`presence-joined`, `presence-left`, and `presence-updated`. Each has
`channelType`, `channelId`, and one `participant` with the same three safe
fields. `presence-updated` covers a second connection for the same user and a
channel-permission change. A `presence-left` participant has `connectionCount:
0`. Clients must apply deltas idempotently because leave/revoke cleanup can race
with lease expiry or reconnection.

The same root connection may receive a user-scoped `resource-event` control
frame after a Core mutation commits. Its payload is a minimal invalidation
hint, not a resource snapshot:

```json
{
  "version": 1,
  "type": "resource-event",
  "eventId": "UUID",
  "eventType": "RootShelfPermissionChanged",
  "resourceId": "UUID",
  "targetUserPublicId": "UUID",
  "change": "permission_updated",
  "permission": "write"
}
```

`eventId` is stable across Kafka redelivery and is the client de-duplication
key. Permission and root-shelf events are sent only to the affected user's
active connections. BlockPack update/delete events are sent to connections
currently subscribed to that BlockPack. Reconnect does not replay historical
resource events; the client refetches canonical REST or GraphQL state.

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

`authenticate` is deliberately rejected with `authentication_managed_by_upgrade`;
root connection authentication is not a channel operation. The subscribe
envelope carries the mandatory `channelTicket`.

## Binary Frames

Binary frames never Base64-encode Yjs data and never use JSON block events. Their header is exactly six bytes, followed by raw bytes:

| Offset | Length | Value |
| --- | --- | --- |
| `0` | 1 byte | protocol version (`1`) |
| `1` | 1 byte | binary type: `1` = `yjs-document`, `2` = `awareness` |
| `2` | 4 bytes | unsigned big-endian `connectorChannelId` |
| `6` | remaining | raw Yjs or awareness payload |

The `connectorChannelId` maps the payload to its subscribed `channelType + channelId`; a public binary frame therefore does not repeat the resource identity. Unknown, unsubscribed, malformed, or unsupported binary frames receive an error JSON frame and are never forwarded.

Gateway 僅將 ready channel 的 Yjs/awareness frame 轉送至已分派的 worker。`yjs-document` 是 mutation，channel ticket 必須具有 `permission: "write"`；read channel 收到 `channel_permission_denied`，且 payload 不會送入 worker。awareness 可由 read/write channel 使用，但由 worker 維護 room-level `y-protocols/awareness` state；它不寫入 `Y.Doc`、snapshot 或 durable update log。payload 是 `encodeAwarenessUpdate` 的 raw encoded update，不包裝 y-websocket protocol header。worker 驗證每個 client ID 僅能被一個 connector channel 宣告，拒絕 malformed、重複或跨 connector 的 client ID，並從 authoritative state 重新編碼後 relay。attach 成功且 document initial state 已準備完成後，worker 先送出 internal `attached`，再補送 authoritative Yjs 與 awareness state；detach、socket close、worker resync 會 broadcast 對應 client ID 的 null-state removal，避免 ghost cursor/presence。`yjs-document` payload 是 `Y.encodeStateAsUpdate` 產生的 raw encoded update。每個 update 先套用至同 room 的 validation Y.Doc 並投影計算完整 Block 數量；整筆 update 只有在 `maximumBlockCount` 內才會套用 authoritative memory Y.Doc、append durable update log 及 broadcast。超額 update 會完整拒絕並重建 validation Y.Doc，不會回溯或污染 authoritative state。Yjs update 可重複套用，因此 sender 收到自己的 relay 不影響正確性。

## Errors And Future Lifecycle

All gateway errors are JSON:

```json
{ "version": 1, "type": "error", "requestId": "sub-1", "connectorChannelId": 1, "code": "channel_not_found", "message": "connectorChannelId is not subscribed on this connection" }
```

Stable error codes are `authentication_managed_by_upgrade`, `binary_channel_not_ready`, `block_pack_quota_exceeded`, `channel_backpressure`, `channel_limit_exceeded`, `channel_not_found`, `channel_permission_denied`, `invalid_acknowledgement`, `invalid_binary_frame`, `invalid_channel_id`, `invalid_channel_ticket`, `invalid_channel_type`, `invalid_connector_channel_id`, `invalid_control_frame`, `permission_revoked`, `resource_unavailable`, `room_admission_unavailable`, `room_connection_limit_exceeded`, `resubscribe_required`, `unsupported_binary_type`, `unsupported_channel_type`, `unsupported_control_type`, `unsupported_message_type`, `unsupported_protocol_version`, and `worker_unavailable`.

`permission_revoked` means the user no longer has the requested channel permission. `resource_unavailable` means the BlockPack, Yjs document, SubShelf, or RootShelf is deleted or otherwise inactive. Both errors remove the logical channel; the client must stop its editor/provider and must not continue sending binary frames on that `connectorChannelId`.

Before upgrading a public socket, Go applies an IP-based upgrade rate limit, its per-process connector cap, and a distributed per-user root-connection cap. A rejected upgrade returns HTTP `429` for the user cap or HTTP `503` for gateway/admission availability; the client must not retry in a tight loop. The distributed cap is represented by Redis TTL leases and is refreshed with the transport heartbeat, so an abnormal process exit is recovered automatically after the lease expires.

Core issues a BlockPack channel ticket containing the short-lived, signed room-admission and document-quota policy snapshot: `roomAdmissionPolicyVersion`, `roomAdmissionEnforcementStrategy` (`reject-new-subscriber`), `maximumSubscribers`, `documentQuotaPolicyVersion`, and `maximumBlockCount`. RealtimeGateway accepts only supported versions and strategy, atomically acquires one shared active-subscriber lease, and passes the verified document quota to Yjs Worker. Read and write subscriptions use the same room capacity. `room_connection_limit_exceeded` means the client should close or unsubscribe another active subscriber before retrying. `block_pack_quota_exceeded` means the complete incoming update was rejected; the client must preserve its draft separately and rebuild from authoritative state before attempting a smaller edit. A successful `unsubscribe`, connection close, permission revocation, and lease expiry all release the slot. Core does not synchronously query active subscriber counts during ownership or plan mutations; its committed lifecycle events flow through the transactional outbox and Kafka to RealtimeGateway, which fans them out through its own Redis Pub/Sub channel to detach matching local channels on every instance.

The gateway caps a connection at 64 active channels. Released IDs are not reused during that connection. Public outbound data uses a bounded queue per `connectorChannelId`, with round-robin scheduling between channels. JSON control frames are always scheduled before binary data. Each channel allows at most 256 queued binary frames and 4 MiB of queued binary payload. Awareness is ephemeral: a queued awareness frame replaces the previous queued awareness frame for that channel. Yjs document updates are never silently dropped or coalesced by Go; if their channel queue is full, the gateway detaches only that channel and sends `channel_backpressure`, requiring a resubscribe/resync while unrelated channels remain active. A failed read or a write that cannot complete within 10 seconds closes the physical socket. Go-to-worker multiplexing uses `YJS_WORKER_URLS`, a comma-separated internal endpoint list. Each URL must target the Yjs worker's `/core/realtime/v1` WebSocket route. Each `blockPackId` maps consistently to one endpoint; each endpoint has one long-lived internal WebSocket and a bounded outbound queue. An unavailable worker or a full internal queue rejects the affected channel payload with `worker_unavailable`.

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
| `39` | remaining | raw Yjs/awareness payload; attach carries its JSON quota policy |

Internal types are `1` `attach`, `2` `detach`, `3` `yjs-document`, `4` `awareness`, `5` `resync-required`, `6` `permission-revoked`, `7` `load-yjs-document`, `8` `yjs-document-loaded`, `9` `append-yjs-update` (legacy single update), `10` `yjs-update-persisted`, `11` `yjs-persistence-failed`, `12` `apply-block-projection`, `13` `block-projection-applied`, `14` `block-projection-failed`, `15` `append-yjs-update-batch`, `16` `load-compactable-yjs-document`, `17` `compactable-yjs-document-loaded`, `18` `apply-compacted-yjs-document`, `19` `yjs-document-compacted`, `20` `yjs-document-compaction-failed`, `21` `attached`, and `22` `block-pack-quota-exceeded`.

`attach` and `detach` are idempotent. Attach carries JSON `{ "version": 1, "maximumBlockCount": N }`; invalid or unsupported policy data is rejected. A first attach asks Core for a binary cold-load payload: `lastUpdateSequence(int64)`, `compactedUntilSequence(int64)`, `projectedUntilSequence(int64)`, snapshot length/state-vector length/update count (`uint32` each), snapshot bytes, state-vector bytes, then ordered update entries of `updateSequence(int64)`, payload length (`uint32`), raw update bytes. The worker registers the subscriber and materializes the document before it sends `attached`; RealtimeGateway starts accepting public binary frames and emits `subscribed` only after receiving that acknowledgement.

During an internal worker reconnect, an already-ready public channel may briefly
produce updates before the worker finishes its cold load. The worker buffers at
most 64 such updates and 256 KiB per room, validates them in arrival order after
materialization, and requests channel resync when that transport-only buffer is
exhausted. These hard limits are not plan quotas and are never placed in tickets.
The binary cold-load state is also capped at 64 MiB before parsing. This is a
memory-safety transport guard, not `maximumDocumentBytes`; introducing a product
document-size quota still requires a dedicated plan/config contract.

`append-yjs-update-batch` carries `[persistenceBatchId:16][originConnectionId:16, zero UUID when mixed][merged raw Yjs update:n]`. Go appends it transactionally and returns an application reply with its authoritative `updateSequence(int64)`; the worker does not wait for that reply on the editor message path. It serializes commands per BlockPack, broadcasts after Kafka's replicated producer ACK, and applies an optimistic local sequence until Core confirms the durable watermark. `persistenceBatchId` makes a retry idempotent when PostgreSQL committed but the ACK was lost. A rejected or timed-out application reply causes a room resync; on an internal worker reconnect, Go replays `attach` for every active channel assigned to that worker before it forwards a client payload. When replay cannot be completed, Go emits `resync_required` to that public channel and waits for the client to resubscribe.

`apply-block-projection` carries UTF-8 JSON `{ schemaId, schemaVersion, projectedSequence, blocks }`; the BlockPack id is the internal frame `channelId`. This request is accepted only over Go-established private worker connections, not through public WebSocket or REST routes. Go validates the schema and durable sequence, bulk applies the BlockTable projection, and returns JSON `{ applied, projectedUntilSequence }` with `block-projection-applied`; malformed, stale-invalid, or failed requests receive `block-projection-failed`.

The internal implementation uses a bounded outbound queue per worker. Queue
exhaustion or a dead worker closes affected logical channels with
`worker_unavailable`; public sockets may remain open for unrelated channels.
The per-channel queue and backpressure rules above are part of this protocol.
