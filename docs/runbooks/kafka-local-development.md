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
YjsWorker command/reply topics and their DLQs:
`notezy.yjsworker.core.command.v1`,
`notezy.yjsworker.core.command.v1.dlq`,
`notezy.core.yjsworker.reply.v1`, and
`notezy.core.yjsworker.reply.v1.dlq`. Automatic topic creation is disabled, so
an absent topic is always a provisioning failure.

## Runtime configuration

`internal/platform/config.Kafka()` reads these values. Core and
RealtimeGateway establish a Kafka client at startup; the other runtimes do not
connect until they own a producer or consumer.

| Variable | Default | Purpose |
| --- | --- | --- |
| `KAFKA_BROKERS` | `127.0.0.1:9094` | Comma-separated broker addresses. |
| `KAFKA_CLIENT_ID` | `notezy-runtime` | Kafka client identity. |
| `KAFKA_CONSUMER_GROUP` | `notezy-runtime` | Default group used by a runtime consumer. |
| `KAFKA_DIAL_TIMEOUT_SECONDS` | `3` | Broker connection timeout. |
| `KAFKA_TLS_ENABLED` | `false` | Enable TLS with system roots or the file values below. |
| `KAFKA_TLS_CA_FILE` | empty | PEM CA file path inside the runtime container. |
| `KAFKA_TLS_CERT_FILE`, `KAFKA_TLS_KEY_FILE` | empty | Optional mTLS pair; both are required together. |
| `KAFKA_TLS_SERVER_NAME` | empty | Optional TLS certificate server name override. |
| `KAFKA_SASL_MECHANISM` | empty | `PLAIN`, `SCRAM-SHA-256`, or `SCRAM-SHA-512`. |
| `KAFKA_SASL_USERNAME`, `KAFKA_SASL_PASSWORD` | empty | Required whenever SASL is configured. |
| `KAFKA_CONSUMER_MAXIMUM_ATTEMPTS` | `3` | Per-delivery handler attempts before a transient error enters DLQ. |
| `KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF_MILLISECONDS` | `250` | First consumer retry delay. |
| `KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF_MILLISECONDS` | `5000` | Cap for consumer exponential retry backoff. |
| `KAFKA_CONSUMER_MAXIMUM_POLL_RECORDS` | `100` | Bound on records held while rebalance is blocked. |

TLS/SASL settings are never committed as Compose secrets. A future deployment
must mount certificate files and inject credentials through its secret manager.

## Health and degraded behavior

Core and RealtimeGateway start even if the initial Kafka ping fails. They log a
degraded-mode warning; their `/healthz` remains process liveness, while
`/readyz` returns `503` until a broker ping succeeds. Every readiness request
performs the next ping, so readiness recovers without restarting the runtime.

Do not treat degraded readiness as permission to lose lifecycle events: Core
persists them in PostgreSQL outbox rows before committing the associated domain
mutation. The relay publishes those rows when Kafka recovers, and
RealtimeGateway consumes them before fan-out to its own Redis Pub/Sub channel.

## Telemetry names

`internal/platform/kafka` owns the common OpenTelemetry metric names. The
outbox relay and consumers must record these instead of inventing runtime-local
variants:

| Metric | Use |
| --- | --- |
| `kafka.broker.ping.duration`, `kafka.broker.ping.count` | Startup/readiness broker reachability. |
| `kafka.publish.duration`, `kafka.publish.count`, `kafka.publish.failure.count` | Producer and outbox relay results. |
| `kafka.consume.duration`, `kafka.consume.count` | Successful consumer handling. |
| `kafka.consumer.lag` | Per topic/group lag observation. |
| `kafka.retry.count`, `kafka.dlq.count` | NOT-33 retry and dead-letter behavior. |
