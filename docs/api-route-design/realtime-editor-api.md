# Realtime Editor API Design

## Purpose

This document defines the backend boundary between the authenticated HTTP API,
the public Realtime Gateway, the Yjs worker, and a BlockPack editor. It owns
the ticket APIs and their relationship to the WebSocket protocol; binary frame
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

All endpoints use the normal authenticated REST pipeline and the standard
response envelope. Routes are rooted at `/api/development/v1`.

| Method | Path | Responsibility |
| --- | --- | --- |
| `POST` | `/realtime/connection/ticket` | Issue a short-lived capability for one root WebSocket connection. |
| `POST` | `/realtime/channel/block-pack/ticket` | Issue a BlockPack-scoped read or write channel capability. |
| `GET` | `/realtime/block-pack/:blockPackId/participants` | Return ephemeral active-presence data for a BlockPack. |

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
metadata, and sequence checkpoints required by the protocol.

### Presence

Participant data is derived from active Redis subscriber leases and is not an
access-control source. An empty response means that no active room connection
was observed; it does not affect membership or permission state.

## Channel lifecycle

1. The Gateway validates the connection ticket during WebSocket upgrade and
   emits `ready` after a successful connection.
2. A `subscribe` control frame carries the channel ticket. The Gateway checks
   that connection and channel claims belong to the same user, authorizes the
   current hierarchy, allocates a connection-local `connectorChannelId`, and
   attaches the channel to the Yjs worker.
3. The worker cold-loads durable Yjs state and sends the complete document
   state through the subscribed channel.
4. On `unsubscribe`, connection close, permission revocation, resource
   unavailability, or resync, the Gateway detaches the channel and releases
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

`permission_revoked` and `resource_unavailable` both terminate the logical
channel. `resync_required` stops live writes until the channel performs a
normal ticket-and-subscribe recovery. The backend never reconstructs an active
document from `BlockTable` rows.

## Related design documents

* [Realtime Protocol](../system-design/realtime-protocol.md)
* [Yjs Collaboration](../system-design/yjs-collaboration.md)
* [RootShelf Sharing API](root-shelf-sharing.md)
