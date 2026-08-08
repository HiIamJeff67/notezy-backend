# Integration test data

Put database seeds, Kafka event fixtures, Redis snapshots, and cross-runtime
contract fixtures in this directory. Runtime-owned E2E data belongs under the
owning runtime's `test/testdata/` directory instead.

## Kafka broker tests

Run `make test-integration-kafka`. When `KAFKA_BROKERS` is set, the tests use
that broker; otherwise they start a temporary Kafka Testcontainer. Set
`KAFKA_BROKERS=localhost:9094` to use the development Compose broker. The tests
require `NOTEZY_RUN_INTEGRATION=1` and are intentionally skipped by ordinary
`make test-integration` runs.
