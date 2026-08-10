# Notification Frontend Integration

This document is the frontend integration contract for Linear `NOT-71`.
The backend Notification Runtime, Gateway routes, Core lifecycle handling,
RealtimeGateway delivery path, and frontend contract are implemented. The
remaining NOT-71 work is live staging verification of production metrics,
dead-letter alerting, and end-to-end delivery.

## Boundaries

The frontend talks only to the public Gateway API and the public
RealtimeGateway WebSocket. It must not call Notification Runtime internal
HTTP endpoints, Core, Kafka, or Notification PostgreSQL directly.

```text
HTTP API:  Browser -> Gateway -> Notification Runtime
WebSocket: Browser -> RealtimeGateway
Events:    Notification Runtime -> Kafka -> RealtimeGateway -> WebSocket
```

The current authenticated user is always taken from the access-token context.
The client must not send a user ID to select another user's notifications.

For browser requests, include cookies (`credentials: "include"` in `fetch`,
or the equivalent TanStack Query client configuration). The Gateway may refresh
the access token during a request; the browser must accept the resulting
`Set-Cookie` headers. Do not copy access or refresh tokens into JSON, query
parameters, or WebSocket URLs.

## HTTP routes

Development API base URL:

```text
/api/development/v1
```

All routes require the normal Gateway authentication cookies. The Gateway
returns the public `ClientResponse` envelope:

```json
{
  "success": true,
  "data": {},
  "exception": null
}
```

### Search private notifications

```http
GET /api/development/v1/notifications/?first=20&after=<endEncodedSearchCursor>
```

Query parameters:

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `first` | integer | no | `1..100`; omit to use the backend default (`50`). |
| `after` | opaque string | no | Use `data.searchPageInfo.endEncodedSearchCursor` from the previous page. |

The response follows the same connection shape used by the GraphQL search
contracts:

```json
{
  "success": true,
  "data": {
    "searchEdges": [
      {
        "encodedSearchCursor": "<opaque-cursor>",
        "node": {
          "id": "4e4b3c2e-2ae4-4c5f-90fd-6e92ef2f4a19",
          "recipientUserPublicId": "00000000-0000-0000-0000-000000000000",
          "type": "news",
          "priority": "normal",
          "templateKey": "news",
          "templateVersion": 1,
          "payload": {
            "title": "Welcome to Notezy",
            "summary": "Your Notezy account is ready.",
            "body": "Start organizing your notes, shelves, and routines in one place."
          },
          "createdAt": "2026-08-09T12:00:00Z",
          "readAt": null,
          "deletedAt": null,
          "expiresAt": null
        }
      }
    ],
    "searchPageInfo": {
      "hasNextPage": true,
      "hasPreviousPage": false,
      "startEncodedSearchCursor": "<opaque-cursor>",
      "endEncodedSearchCursor": "<opaque-cursor>"
    },
    "totalCount": 1,
    "searchTime": 1.42
  },
  "exception": null
}
```

The cursor is opaque and must not be decoded or constructed by the frontend.
When placing it in the URL, let the HTTP client URL-encode it (for example,
use `URLSearchParams` rather than string concatenation).
When `searchPageInfo.hasNextPage` is `false`, there is no next page. The
backend does not return deleted or expired notifications from this endpoint.
The query uses `(createdAt, id)` as a stable keyset, so equal timestamps do
not cause duplicate or skipped notifications.

### Unread count

```http
GET /api/development/v1/notifications/unread-count
```

Response:

```json
{
  "success": true,
  "data": { "count": 3 },
  "exception": null
}
```

### Mark notifications as read

```http
PATCH /api/development/v1/notifications/read
Content-Type: application/json

{
  "notificationIds": [
    "4e4b3c2e-2ae4-4c5f-90fd-6e92ef2f4a19"
  ]
}
```

The request is scoped to the authenticated user. IDs belonging to another
user are not modified.

Response:

```json
{
  "success": true,
  "data": { "updatedCount": 1 },
  "exception": null
}
```

### Soft-delete notifications

```http
DELETE /api/development/v1/notifications/
Content-Type: application/json

{
  "notificationIds": [
    "4e4b3c2e-2ae4-4c5f-90fd-6e92ef2f4a19"
  ]
}
```

Response:

```json
{
  "success": true,
  "data": { "deletedCount": 1 },
  "exception": null
}
```

This is a soft delete. The client should remove the item from the visible
inbox immediately and refetch if it needs server confirmation.

## Realtime ticket APIs

Notification delivery uses the public RealtimeGateway WebSocket. The ticket
APIs are served by the Gateway and are authenticated with the same cookies as
the notification HTTP routes.

### Connection ticket

```http
POST /api/development/v1/realtime/connection/ticket
```

The request has no JSON body. The browser's `User-Agent` header is captured by
the Gateway and bound to the short-lived ticket. The successful response is a
`ClientResponse` whose `data` contains:

```json
{
  "realtimeEndpoint": "/realtime/development/v1",
  "realtimeProtocolVersion": 1,
  "connectionTicket": "<short-lived-ticket>",
  "expiresAt": "2026-08-09T12:05:00Z"
}
```

Use `data.realtimeEndpoint` with the public RealtimeGateway origin. The
connection ticket is a single-use capability and must be sent as the only
WebSocket subprotocol. It is not a bearer token for REST calls.

### BlockPack channel ticket

Issue one channel ticket for each BlockPack subscription:

```http
POST /api/development/v1/realtime/channel/block-pack/ticket
Content-Type: application/json

{
  "blockPackId": "4e4b3c2e-2ae4-4c5f-90fd-6e92ef2f4a19",
  "permission": "read"
}
```

`permission` is `read` or `write`; Core checks the authenticated user's
current BlockPack permission before issuing the ticket. The response includes
`channelTicket`, `channelId`, `permission`, room/document metadata, sequence
checkpoints, and `expiresAt`. The channel ticket is also single-use and is
sent in the WebSocket `subscribe` frame, not in the URL.

## Error handling

Errors use the same public envelope:

```json
{
  "success": false,
  "data": null,
  "exception": {
    "reason": "Unauthorized",
    "domain": "Auth",
    "operation": "JWT",
    "message": "Authentication is required.",
    "retryable": false
  }
}
```

The frontend should treat `retryable: true` as eligible for a bounded retry,
but should not retry validation, authentication, or permission errors. Unknown
notification types must use a fallback renderer rather than crashing the
notification panel.

## WebSocket delivery

RealtimeGateway development WebSocket URL:

```text
wss://<host>/realtime/development/v1
```

The normal flow is:

1. Call `POST /api/development/v1/realtime/connection/ticket` through the
   Gateway.
2. Connect directly to RealtimeGateway using `data.realtimeEndpoint` and the
   returned connection ticket as the single `Sec-WebSocket-Protocol` value.
3. Wait for the `ready` frame. A BlockPack `subscribe` frame is optional and
   is independent of notification delivery.
4. Listen on the root connection for frames whose `type` is `notification`.

Example browser connection:

```ts
const ticketHttpResponse = await fetch(
  "/api/development/v1/realtime/connection/ticket",
  { method: "POST", credentials: "include" },
);
const ticketEnvelope = await ticketHttpResponse.json();
const socket = new WebSocket(
  `${realtimeGatewayOrigin}${ticketEnvelope.data.realtimeEndpoint}`,
  [ticketEnvelope.data.connectionTicket],
);
```

The WebSocket upgrade does not use access-token cookies or an Authorization
header. RealtimeGateway authenticates the upgrade solely with the signed,
single-use connection ticket and the browser `User-Agent` binding.

Notification frame shape:

```json
{
  "version": 1,
  "type": "notification",
  "eventId": "8dc5f1f4-a5ba-4c7d-91b2-4f82cb9ef4dd",
  "notificationId": "4e4b3c2e-2ae4-4c5f-90fd-6e92ef2f4a19",
  "notificationType": "news",
  "priority": "normal",
  "templateKey": "news",
  "templateVersion": 1,
  "payload": {
    "title": "Welcome to Notezy",
    "summary": "Your Notezy account is ready.",
    "body": "Start organizing your notes, shelves, and routines in one place."
  },
  "createdAt": "2026-08-09T12:00:00Z",
  "expiresAt": null
}
```

The WebSocket frame is a live-delivery hint, not the history source of truth.
On reconnect, the frontend must call the list endpoint with its last known
cursor and deduplicate by `notificationId` (or `eventId`). A frame may arrive
before or after the HTTP response, so the UI store must merge both sources
idempotently.

Notification frames are delivered on every authenticated root connection; the
client does not need to subscribe to a BlockPack to receive them. The frame's
`payload` is JSON with a type-specific shape selected by `templateKey` and
`templateVersion`.

## Suggested Zod contracts

The frontend can keep the dynamic payload boundary explicit while validating
the stable envelope:

```ts
const notificationSchema = z.object({
  id: z.string().uuid(),
  recipientUserPublicId: z.string().uuid(),
  type: z.enum(["news", "warning", "important"]),
  priority: z.enum(["low", "normal", "high", "critical"]),
  templateKey: z.string(),
  templateVersion: z.number().int().positive(),
  payload: z.record(z.string(), z.unknown()),
  createdAt: z.string().datetime(),
  readAt: z.string().datetime().nullable(),
  deletedAt: z.string().datetime().nullable(),
  expiresAt: z.string().datetime().nullable(),
});
```

Use `templateKey` and `templateVersion` to select a type-specific payload
schema. Keep an unknown-template fallback so a newer backend contract does not
break an older frontend deployment.

## TanStack Query integration checklist

- Keep `notifications` and `unread-count` as separate query keys.
- Use `searchPageInfo.endEncodedSearchCursor` as the cursor for infinite pagination.
- Insert a WebSocket notification into the list cache only if its
  `notificationId` is not already present.
- Increment or invalidate unread count for a new unread frame.
- Optimistically mark read and roll back on a failed mutation.
- Optimistically remove soft-deleted items and invalidate the list after the
  mutation settles.
- Reconnect WebSocket, then refetch the first page to close any delivery gap.
- Render `news`, `warning`, and `important` separately; render unknown types
  with a neutral fallback card.

## Backend status for NOT-71

Already available:

- Notification PostgreSQL runtime and JSONB payload persistence.
- Core outbox producer, including registration welcome notification.
- Notification consumer, idempotency inbox, and Notification outbox relay.
- RealtimeGateway Kafka consumer and live WebSocket notification frames.
- Gateway authenticated private-search, unread count, read, and soft-delete routes.
- Kafka contract coverage for Core notification requests and user-deletion
  lifecycle events (run with `NOTEZY_RUN_INTEGRATION=1`).

Still tracked in NOT-71:

- Production metrics/DLQ deployment verification in the staging environment.
- End-to-end acceptance against the deployed Gateway, Notification Runtime,
  Kafka, and RealtimeGateway runtimes.
