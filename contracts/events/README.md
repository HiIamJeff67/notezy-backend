# Shared event transport contract

This package contains only the runtime-neutral Kafka envelope and primitive
transport values used by `internal/platform/kafka`:

- `EventEnvelope[D]`
- `Topic`
- `EventType`
- `AggregateType`
- `TraceMetadata`

It must not contain a business event payload, topic name, consumer group, or
runtime-specific operation. Those contracts belong to the runtime that owns
the interaction, such as `contracts/core/v1`, `contracts/durablejob/v1`, or
`contracts/yjsworker/v1`.
