# Kafka Event Contracts

## Scope

This document defines the versioned Core lifecycle events used across service
boundaries. It is the source of truth for the event contract in
[`contracts/core/v1/events`](../../contracts/core/v1/events/). Core writes these events
through its transactional outbox and relay. RealtimeGateway consumes them,
fans them out through its runtime-owned Redis Pub/Sub channel, and executes
local socket/channel cleanup. NOT-33 supplies consumer idempotency, retry, and
DLQ behavior.

Events state a fact already committed by the producer and a policy already
decided by Core. Consumers execute that fact or policy; they never derive
membership, quota, hierarchy, or room admission rules from local state.

## Envelope

Every message is a JSON `eventscontract.EventEnvelope[D]` with these fields.

| Field | Meaning |
| --- | --- |
| `schemaVersion` | Contract major version, initially `v1`. |
| `eventId` | Globally unique UUID used by consumer idempotency. |
| `eventType` | Stable lifecycle fact name. |
| `aggregateType`, `aggregateId` | Owner of ordering and causal identity. |
| `kafkaKey` | Exact `aggregateId` UUID string. |
| `occurredAt` | UTC time at which Core committed the fact. |
| `correlationId` | Required workflow/request identity propagated across related work. |
| `causationId` | Optional parent event UUID. |
| `trace` | Optional W3C `traceParent` and `traceState`; no browser credentials. |
| `data` | Event-type-specific minimal payload. |

`User` aggregate IDs and `TargetUserPublicId` values are user public UUIDs.
Core internal user IDs must never cross this boundary.

## Initial lifecycle events

| Event type | Aggregate | Data | RealtimeGateway action |
| --- | --- | --- |
| `BlockPackAccessRevoked` | one `BlockPack` | optional `targetUserPublicId` | Detach the target user's channel, or every channel when absent. |
| `BlockPackRoomPolicyChanged` | one `BlockPack` | policy version, strategy, maximum subscribers | Apply the attested policy to later admissions; the current v1 strategy rejects new subscribers. |
| `RootShelfPermissionRevoked` | one `RootShelf` | optional `targetUserPublicId` | Lifecycle/audit fact. Core also emits one BlockPack revoke per affected descendant for channel detachment. |
| `UserSessionsRevoked` | one `User` public UUID | empty | Fan out a local session revoke and close the user's active RealtimeGateway connections. |

`SubShelf` is an initial aggregate type for future lifecycle facts. It has no
separate v1 event yet: mutations that must detach active rooms emit the affected
one-BlockPack `BlockPackAccessRevoked` events.

## Topic and ordering strategy

The initial topic is `notezy.core.lifecycle.v1`.

| Property | Rule |
| --- | --- |
| Producer | Core only. |
| Consumer group | A runtime-specific group, for example `notezy-realtimegateway-lifecycle-v1`. |
| Partition key | `aggregateId.String()` exactly; it is copied into `kafkaKey`. |
| Ordering | Kafka ordering is guaranteed only for the same aggregate key/partition. No cross-aggregate ordering is assumed. |
| Partitions | Deploy with 12 production partitions initially; increase only by adding partitions and never rely on cross-key order. |
| Retention | Seven days, delete-based cleanup; this topic is a transport log, not the event-sourced source of truth. |
| Replication | Production uses replication factor three and `min.insync.replicas=2`. |

A bulk lifecycle mutation performs a batch outbox insert. It emits one row and
one Kafka record for each affected BlockPack, keyed by that BlockPack ID; it
never emits a list of unrelated BlockPack IDs in one record. This preserves the
only ordering guarantee consumers may rely on.

## Compatibility and payload limits

Within a major contract version, producers may add optional JSON fields only.
Consumers must ignore unknown fields. Renaming, changing a field type or
meaning, making an optional field required, or changing an event's ordering
meaning requires `contracts/events/v2` and a new `.v2` topic. Producers publish
the new topic independently until all consumers migrate.

Lifecycle payloads contain identifiers and already-decided policy values only.
They must not contain complete entity snapshots, names/profile fields, database
internal IDs, browser cookies, access/refresh/delegation tokens, Yjs updates,
awareness, presence, Redis leases, or WebSocket frames.

## YjsWorker command/reply

Yjs persistence is a separate, versioned Kafka boundary in
[`contracts/yjsworker/v1`](../../contracts/yjsworker/v1/). YjsWorker produces
an `EventEnvelope` carrying a `CommandEnvelope` to
`notezy.yjsworker.core.command.v1`, keyed by `blockPackId`. Core consumes it
with the `notezy-core-yjsworker-v1` group and writes a matching reply envelope
to its transactional outbox. The outbox relay publishes the reply to
`notezy.core.yjsworker.reply.v1`, again keyed by the same BlockPack ID.

The command family is `LoadYjsDocument`, `AppendYjsUpdate`,
`LoadCompactableYjsDocument`, `ApplyCompactedYjsDocument`, and
`ApplyBlockProjection`. Each carries its own command UUID, correlation ID,
optional causation ID, trace metadata, producer identity, and BlockPack
partition key. Core echoes those values in its reply. Append deduplicates by
`persistenceBatchId`; compaction and projection use their durable sequence
checkpoints, so Kafka redelivery cannot apply a newer state twice.

Every Core mutation and its reply outbox row share one database transaction.
That prevents a committed mutation from losing its reply if Core stops before
Kafka is available. YjsWorker treats reply timeout as retryable and reuses its
persistence batch ID; invalid envelopes are routed to the normal command DLQ.
RealtimeGateway does not relay Yjs persistence or projection HTTP requests.

## Consumer reliability and idempotency

Every runtime consumer is built on `internal/platform/kafka.Consumer`. It uses
manual offset commits, bounded `PollRecords`, and Kafka's
`BlockRebalanceOnPoll`: an offset is committed only after its handler has
completed successfully, or after the original record has been durably written
to its dead-letter topic. The consumer always releases the rebalance gate when
the batch ends. A process crash, handler panic, failed dead-letter publish, or
failed offset commit therefore leaves the source record uncommitted for
redelivery.

Kafka and the Core outbox are at-least-once. The consumer platform cannot
atomically mark an arbitrary runtime side effect as complete, so it must not
pretend to own a global deduplication table. Each concrete handler owns its
idempotency boundary: it uses `eventId` (or the event's monotonic aggregate
version when a future contract defines one) in the same transaction or
idempotent store operation as its side effect. A redelivered event must then be
a no-op and still return success. RealtimeGateway's Redis relay may redeliver,
but local detach closes an existing socket or removes an existing channel; a
repeated revoke finds neither and is therefore a no-op. NOT-35 applies this rule to RealtimeGateway;
future Core, DurableJob, and Email consumers apply it in their own data
boundaries.

Handlers return `*kafka.ConsumerError` to classify a failure:

| Classification | Handling |
| --- | --- |
| `Transient` (the default for ordinary errors) | Retry locally with bounded exponential backoff. After `KAFKA_CONSUMER_MAXIMUM_ATTEMPTS`, write to DLQ. |
| `PoisonMessage` | Write to DLQ immediately; retrying cannot make the handler accept it. |
| `SchemaIncompatible` | Write to DLQ immediately. Envelope decode/contract validation produces this classification automatically. |

The retry limit is intentionally per active delivery lease. A crash or
rebalance before a terminal commit restarts delivery, preserving data rather
than silently discarding it. DLQ delivery contains the original topic,
partition, offset, key, unmodified source bytes, event ID when decodable,
classification, attempt count, and failure timestamp.

The initial dead-letter topic is `notezy.core.lifecycle.v1.dlq`, configured
with 30-day delete retention. After the cause has been corrected, operators
inspect the DLQ record, verify its original envelope is now compatible, and
re-publish its preserved `key` and `value` to `sourceTopic`. This re-drive goes
through the normal consumer/idempotency path; no consumer offsets are edited
manually. Poison or incompatible records should be corrected at the producer
or contract layer before re-drive.

Consumer telemetry uses `kafka.consume.*`, `kafka.consumer.lag`,
`kafka.retry.count`, and `kafka.dlq.count`, alongside structured failure logs.
`traceParent`/`traceState` in the event envelope are extracted before the
handler runs, so handler spans continue the producer trace.
