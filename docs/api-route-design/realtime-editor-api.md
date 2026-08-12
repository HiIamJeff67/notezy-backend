# Realtime Editor API Design

## Purpose

This document defines the backend boundary between the authenticated HTTP API,
the public Realtime Gateway, the Yjs worker, and a BlockPack editor. It owns
the ticket APIs, direct RealtimeGateway presence query, and their relationship to the WebSocket protocol; binary frame
formats and durable Yjs semantics are defined in
[Realtime Protocol](../system-design/realtime-protocol.md) and
[Yjs Collaboration](../system-design/yjs-collaboration.md).

## Source-of-truth boundaries

| Data | Authoritative use | Not responsible for |
| --- | --- | --- |
| Active `Y.Doc` | collaborative editing and synchronization | REST block mutation |
| `BlockPackYjsDocument` plus update tail | durable document state | direct public access |
| `BlockTable` | REST/GraphQL projection read and search | initializing or repairing a live document |

The backend must never accept the same editor mutation through both REST and
Yjs. There is no public REST Block mutation API.

## HTTP API

Ticket endpoints use the authenticated API pipeline and are rooted at
`/api/development/v1`. Presence is served directly by RealtimeGateway under
`/realtime/development/v1`.

| Method | Path | Responsibility |
| --- | --- | --- |
| `POST` | `/realtime/connection/ticket` | Issue a short-lived capability for one root WebSocket connection. |
| `POST` | `/realtime/channel/block-pack/ticket` | Issue a BlockPack-scoped read or write channel capability. |
| `GET` | `/realtime/block-pack/:blockPackId/participants` | Return RealtimeGateway-owned ephemeral active presence directly from its Redis lease store. User profile details remain a separate Core API concern. |

### Connection ticket

The connection-ticket response contains `realtimeEndpoint`,
`realtimeProtocolVersion`, `connectionTicket`, and `expiresAt`. The ticket is
valid only for the WebSocket upgrade and is sent as the sole
`Sec-WebSocket-Protocol` value. It identifies the authenticated user but does
not authorize any BlockPack.

### BlockPack channel ticket

The request body contains a `blockPackId` and a requested `permission`
(`read` or `write`). The backend validates the active
`RootShelf -> SubShelf -> BlockPack -> BlockPackYjsDocument` hierarchy and the
caller's effective RootShelf permission before issuing a ticket.

`Read`, `Write`, `Admin`, and `Owner` may request `read`; only `Write`,
`Admin`, and `Owner` may request `write`. The response contains the signed
ticket, expiry, channel identity, effective permission, room/document schema
metadata, `documentQuotaPolicyVersion`, `maximumBlockCount`, and the sequence
checkpoints required by the protocol. The signed ticket carries the same BlockPack
quota policy plus the private room-admission policy snapshot (version, strategy,
and maximum subscribers). The public quota fields are client UX hints; the signed
claims verified by RealtimeGateway and enforced by Yjs Worker are authoritative.

### Presence

Participant data is derived from active Redis subscriber leases and is not an
access-control source. An empty response means that no active room connection
was observed; it does not affect membership or permission state.

## Channel lifecycle

1. RealtimeGateway validates the connection ticket during WebSocket upgrade and
   emits `ready` after a successful connection.
2. A `subscribe` control frame carries the channel ticket. RealtimeGateway checks
   that connection and channel claims belong to the same user, applies the
   ticket's signed room-admission policy, allocates a connection-local
   `connectorChannelId`, and attaches the channel to the Yjs worker.
3. The worker cold-loads durable Yjs state, registers the subscriber, and returns
   an internal `attached` acknowledgement. RealtimeGateway emits public
   `subscribed` only after that acknowledgement, followed by the complete document
   state.
4. On `unsubscribe`, connection close, permission revocation, resource
   unavailability, or resync, RealtimeGateway detaches the channel and releases
   all associated leases.

`connectorChannelId` is valid only for one root connection and must never be
treated as a BlockPack identifier or persisted as resource identity.

## Projection reads

Block and BlockPack REST/GraphQL reads expose the `BlockTable` materialized
projection. BlockPack responses provide `lastUpdateSequence`,
`compactedUntilSequence`, `projectedUntilSequence`, and
`isProjectionCurrent`; a stale projection is an observation signal only and
does not change the durable Yjs source of truth.

## Failure semantics

`permission_revoked`, `resource_unavailable`, and `block_pack_quota_exceeded`
terminate the logical channel. `resync_required` stops live writes until the
channel performs a normal ticket-and-subscribe recovery. A quota rejection does
not mutate, persist, project, or broadcast the rejected update. The client must
preserve the rejected draft separately, rebuild a clean Y.Doc from authoritative
state, and let the user reduce the content before retrying. The backend never
reconstructs an active document from `BlockTable` rows.

## Related design documents

* [Realtime Protocol](../system-design/realtime-protocol.md)
* [Yjs Collaboration](../system-design/yjs-collaboration.md)
* [RootShelf Sharing API](root-shelf-sharing.md)
