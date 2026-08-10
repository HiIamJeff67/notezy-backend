# Kafka topic bootstrap

Kafka topic names and their broker settings are maintained in
`shared/platform/kafka/topics`. That package exposes `topics.TopicSpec` and
`topics.All()`. The platform `TopicProvisioner` accepts those specs and
idempotently creates missing topics and dead-letter topics while leaving
existing topics unchanged.

Each queue has its own `...TopicSpec()` factory. The factory explicitly sets
partitions, replication factor, retention, cleanup policy, minimum in-sync
replicas, and dead-letter retention when applicable. `TopicProvisioner` only
validates and applies those values; it does not fill in broker defaults. This
keeps queue policy owned by the runtime contract that defines the queue and
makes an incomplete topic specification fail during bootstrap.

The development Compose stack runs the same implementation through the CLI
service `notezy-kafka-init`:

```sh
make kafka-topics
```

The command reads `KAFKA_BROKERS`, `KAFKA_DIAL_TIMEOUT`, and the Kafka security
settings from the environment. Runtime services depend on the init service
completing successfully before they start consuming or producing events.

Integration tests use `infra/docker/docker-compose.integration.yaml` to start the same
Kafka image and then call the same provisioner after the broker is healthy.
This keeps local and CI broker startup consistent without duplicating a second
topic list in test code or shell scripts.

`topics.TopicSpec` is deliberately a creation specification, not a runtime topic
client. It describes partitions, replication factor, retention, cleanup policy,
minimum in-sync replicas, and whether a dead-letter topic should be created.
It also carries the explicit dead-letter retention when dead-letter creation is
enabled. Changing a specification does not mutate an existing topic; broker
migrations must be explicit operational changes.
