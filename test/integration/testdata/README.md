# Integration test data

Put database seeds, Kafka event fixtures, Redis snapshots, and cross-runtime
contract fixtures in this directory. Runtime-owned E2E data belongs under the
owning runtime's `test/testdata/` directory instead.

## Kafka broker tests

Run `make test-integration-kafka` after the integration Compose stack is
available. The tests use `KAFKA_BROKERS` (defaulting to the integration broker
at `127.0.0.1:19094`) and do not create or manage a Kafka container themselves.
They require `NOTEGIC_RUN_INTEGRATION=1` and are intentionally kept separate from
the ordinary `make test-integration` suite.
