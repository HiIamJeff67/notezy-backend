# Kafka Local Development

## Scope

Kafka is the Phase 3 transport for versioned lifecycle events. Local Compose
runs one KRaft broker and provisions the v1 lifecycle topic. This is a local
and integration environment only; production Kafka topology, credentials, and
delivery automation remain outside this runbook.

The contract and production topic strategy are defined in
[Kafka Event Contracts](../system-design/kafka-event-contracts.md). The local
broker deliberately uses three partitions and replication factor one; the
production strategy is twelve partitions, replication factor three, and
`min.insync.replicas=2`.

## Start and inspect

Start the broker and one-time topic provisioner:

```bash
docker compose up -d notezy-kafka notezy-kafka-init
docker compose logs --follow notezy-kafka-init
```

Inspect a provisioned topic:

```bash
docker compose exec notezy-kafka \
  /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --describe \
  --topic notezy.core.lifecycle.v1
```

Containers connect to `notezy-kafka:9092`. A process started on the host uses
`127.0.0.1:9094` by default. The initial lifecycle topic has delete cleanup,
seven-day retention, three local partitions, and replication factor one. Its
paired dead-letter topic, `notezy.core.lifecycle.v1.dlq`, has the same local
partition count and 30-day retention. The same provisioner also creates the
YjsWorker command/reply and DurableJob maintenance topics:
`notezy.yjsworker.core.command.v1`,
`notezy.yjsworker.core.command.v1.dlq`,
`notezy.core.yjsworker.reply.v1`,
`notezy.core.yjsworker.reply.v1.dlq`,
`notezy.core.durablejob.yjs-maintenance-hint.v1`,
`notezy.durablejob.core.yjs-maintenance-request.v1`,
`notezy.core.yjsworker.maintenance-command.v1`,
`notezy.yjsworker.core.maintenance-result.v1`, and
`notezy.core.durablejob.yjs-maintenance-result.v1` (and their local DLQs).
Automatic topic creation is disabled, so
an absent topic is always a provisioning failure.

## Runtime configuration

`shared/platform/kafka/config.go` reads the infrastructure connection values.
Each runtime keeps its own consumer policy in its local config package. Core,
DurableJob, and RealtimeGateway establish Kafka clients at startup when they
own a producer or consumer.

| Variable | Default | Purpose |
| --- | --- | --- |
| `KAFKA_BROKERS` | `127.0.0.1:9094` | Comma-separated broker addresses. |
| `KAFKA_DIAL_TIMEOUT` | `3s` | Broker connection timeout as a Go duration. |
| `KAFKA_TLS_ENABLED` | `false` | Enable TLS with system roots or the file values below. |
| `KAFKA_TLS_CA_FILE` | empty | PEM CA file path inside the runtime container. |
| `KAFKA_TLS_CERT_FILE`, `KAFKA_TLS_KEY_FILE` | empty | Optional mTLS pair; both are required together. |
| `KAFKA_TLS_SERVER_NAME` | empty | Optional TLS certificate server name override. |
| `KAFKA_SASL_MECHANISM` | empty | `PLAIN`, `SCRAM-SHA-256`, or `SCRAM-SHA-512`. |
| `KAFKA_SASL_USERNAME`, `KAFKA_SASL_PASSWORD` | empty | Required whenever SASL is configured. |
| `KAFKA_CONSUMER_MAXIMUM_ATTEMPTS` | `3` | Per-delivery handler attempts before a transient error enters DLQ. |
| `KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF` | `250ms` | First consumer retry delay as a Go duration. |
| `KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF` | `5s` | Cap for consumer exponential retry backoff as a Go duration. |
| `KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS` | `100` | Bound on records held while rebalance is blocked. |

TLS/SASL settings are never committed as Compose secrets. A future deployment
must mount certificate files and inject credentials through its secret manager.

## Health and degraded behavior

Core and RealtimeGateway start even if the initial Kafka ping fails. They log a
degraded-mode warning; their `/startedz` reports that the process has started,
while `/healthz` returns `503` until the runtime can accept normal operations.
The current startup status is evaluated once; restart the runtime after the
required broker connection becomes available so `/healthz` can report healthy.

Do not treat degraded health as permission to lose lifecycle events: Core
persists them in PostgreSQL outbox rows before committing the associated domain
mutation. The relay publishes those rows when Kafka recovers, and
RealtimeGateway consumes them before fan-out to its own Redis Pub/Sub channel.

## Telemetry names

`shared/platform/kafka` owns the common OpenTelemetry metric names. The
outbox relay and consumers must record these instead of inventing runtime-local
variants:

| Metric | Use |
| --- | --- |
| `kafka.broker.ping.duration`, `kafka.broker.ping.count` | Startup/health broker reachability. |
| `kafka.publish.duration`, `kafka.publish.count`, `kafka.publish.failure.count` | Producer and outbox relay results. |
| `kafka.consume.duration`, `kafka.consume.count` | Successful consumer handling. |
| `kafka.consumer.lag` | Per topic/group lag observation. |
| `kafka.retry.count`, `kafka.dlq.count` | NOT-33 retry and dead-letter behavior. |
