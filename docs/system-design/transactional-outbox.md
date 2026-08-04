# Core transactional outbox

Core persists lifecycle events in `OutboxEventTable` in the same PostgreSQL
transaction as the originating domain mutation. The table is Core-owned and is
the durable handoff boundary between PostgreSQL and Kafka; Kafka is not part of
the PostgreSQL transaction.

```text
Core domain transaction
  ├─ mutate Core-owned tables
  └─ INSERT one or more OutboxEvent rows
       └─ COMMIT

Core outbox relay
  ├─ claim available rows with FOR UPDATE SKIP LOCKED
  ├─ publish the event-contract envelope to Kafka
  ├─ Kafka acknowledgement → mark published
  └─ error → increment publish count and defer availableAt with backoff
```

## Delivery semantics

The relay stores a temporary worker claim. A crashed worker leaves a claim that
expires after `OUTBOX_RELAY_CLAIM_TIMEOUT_SECONDS`; another relay can then claim
and publish the row. If Kafka acknowledges a record but the process crashes
before `publishedAt` is stored, the row is deliberately published again after
the claim expires. Delivery is therefore at-least-once, not cross-system
exactly-once. Consumers must be idempotent; that belongs to NOT-33.

Claiming, marking published rows, scheduling failures, and cleanup are each
set-based database operations. Multiple Core relay processes can run safely:
`FOR UPDATE SKIP LOCKED` prevents them from claiming the same available row at
the same time.

## Contract mapping

`payload` stores the event `Data`; `metadata` stores the remaining versioned
event envelope metadata. The relay reconstructs and publishes
`eventscontract.EventEnvelope[json.RawMessage]`, preserving the schema version,
event/aggregate IDs, Kafka key, correlation/causation identifiers, timestamp,
and trace metadata. The key must equal `aggregateId.String()`.

Domain workflows use `outbox.EnqueueMany(tx, topic, envelopes)` before their
transaction commits. Calling it outside a transaction is rejected by the
repository. Lifecycle producer adoption and the RealtimeGateway consumer are
deliberately owned by NOT-35.

## Configuration and maintenance

| Variable | Default | Purpose |
| --- | --- | --- |
| `OUTBOX_RELAY_BATCH_SIZE` | `100` | Maximum rows claimed per poll. |
| `OUTBOX_RELAY_POLL_INTERVAL_MILLISECONDS` | `1000` | Relay polling interval. |
| `OUTBOX_RELAY_CLAIM_TIMEOUT_SECONDS` | `30` | Crash-recovery claim expiry. |
| `OUTBOX_RELAY_INITIAL_BACKOFF_MILLISECONDS` | `1000` | First retry delay. |
| `OUTBOX_RELAY_MAXIMUM_BACKOFF_SECONDS` | `60` | Maximum exponential retry delay. |
| `OUTBOX_RELAY_RETENTION_HOURS` | `168` | Published-event retention. |
| `OUTBOX_RELAY_CLEANUP_INTERVAL_MINUTES` | `60` | Published-row cleanup cadence. |

The relay emits `outbox.relay.claimed.count`, `outbox.relay.published.count`,
`outbox.relay.failure.count`, `outbox.relay.retry.count`, and
`outbox.cleanup.deleted.count`, and creates a trace span for every relay batch.
Kafka publish telemetry remains under the `kafka.*` metrics defined in the
[Kafka local development runbook](../runbooks/kafka-local-development.md).
