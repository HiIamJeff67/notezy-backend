# Load and soak tests

These scenarios use [Grafana k6](https://grafana.com/docs/k6/latest/) and are
kept outside the Go modules. They are intentionally configurable because the
WebSocket ticket and metrics endpoints depend on the environment under test.

Examples:

```sh
make test-load-websocket \
  REALTIME_GATEWAY_WS_URL='ws://127.0.0.1:7779/realtime/channel/block-pack' \
  REALTIME_CHANNEL_TICKET='<short-lived-ticket>' \
  K6_VUS=20 K6_DURATION=5m

make test-soak-websocket \
  REALTIME_GATEWAY_WS_URL='ws://127.0.0.1:7779/realtime/channel/block-pack' \
  REALTIME_CHANNEL_TICKET='<short-lived-ticket>' \
  K6_VUS=100 K6_DURATION=2h

make test-load-kafka-lag \
  KAFKA_METRICS_URL='http://127.0.0.1:7778/metrics' \
  KAFKA_LAG_THRESHOLD=1000 K6_DURATION=30m
```

The WebSocket scenario measures connection establishment and long-lived room
stability. Kafka lag is read from the runtime metrics endpoint rather than
connecting to Kafka directly from k6; Kafka protocol producers and consumers
are covered by the Go integration suite and NOT-66.
