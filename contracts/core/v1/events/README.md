# Event contracts v1

This package is Core's versioned domain-event boundary for Kafka. It has no
dependency on Core implementation, RealtimeGateway, persistence, or a Kafka
client. The generic envelope and primitive transport types are imported from
`contracts/types/events/`; this package owns only Core lifecycle topics and payloads.

`notezy.core.lifecycle.v1` carries Core lifecycle facts. Every event uses its
aggregate UUID string as `kafkaKey`, so one aggregate remains ordered within a
Kafka partition. `BlockPack` events always describe exactly one BlockPack;
bulk mutations create one outbox row and one event per affected BlockPack.

Resource events on the same topic are the user-scoped UI/cache invalidation
boundary. `RootShelfPermissionChanged`, `RootShelfPermissionRevoked`, and
`RootShelfDeleted` carry the affected user's public UUID. `BlockPackChanged` and
`BlockPackDeleted` omit the target user and are delivered only to connections
currently subscribed to that BlockPack. They contain identifiers and the
already-decided change only; clients must refetch through the normal API when
they need the new resource state.

`TargetUserPublicId` is optional. A present value targets only that user's
RealtimeGateway connections; an omitted value applies to all active channels
for a BlockPack aggregate. It always contains the public user UUID, never
Core's internal user ID. RootShelf permission and deletion events are emitted
once per affected user, so they always carry a target in the resource-event
path.

The envelope and payloads exclude entity snapshots, browser credentials,
cookies, tokens, presence, awareness, and Yjs updates. Producers may add only
backward-compatible optional fields within `v1`; a breaking change requires a
new contract and topic major version.

Core also publishes `YjsMaintenanceHintData` on the Core-to-DurableJob hint
topic. The hint contains only Core's current maintenance metadata; it never
contains a snapshot, state vector, or Yjs update payload.

Core publishes `RoutineTaskCompletedData` on the lifecycle topic after a
prepared DurableJob result has been applied. It contains only task identity,
attempt, worker, and completion metadata; it never contains Core schemas or
database payloads.
