# Presence and event rules

`GET /realtime/development/v1/block-pack/{blockPackId}/participants` returns an ephemeral Redis lease snapshot. It is not an authorization source.

Each participant contains only public user ID, read/write channel permission, and active connection count. An empty list means no active lease was observed. Profile data is intentionally excluded.

After subscription, the socket may emit `presence-joined`, `presence-left`, and `presence-updated`. Apply these deltas idempotently. A left participant has connectionCount zero.

`resource-event` is an invalidation hint, not a resource snapshot. Deduplicate with eventId and refetch canonical REST/GraphQL state. Historical events are not replayed after reconnect. User notifications may also arrive on the root connection and must be treated as transient delivery.

`routine-task-lifecycle` is a user-targeted transient execution hint. It reports `running` when DurableJob begins a RoutineTask handler and `completed` only after Core commits the corresponding result. Deduplicate with eventId; after reconnect or when durable state matters, refetch the canonical RoutineTask and RoutineTaskRecord through the normal API.
