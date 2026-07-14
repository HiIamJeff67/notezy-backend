# Yjs Observability

Development starts Grafana with the LGTM stack. The `Yjs Collaboration Overview` dashboard shows active rooms and subscribers, Go gateway connectors, worker heap and event-loop delay, Yjs operation throughput, p95 latency, and payload throughput.

The Go API exports OTLP over gRPC to `OTEL_EXPORTER_OTLP_GRPC_ENDPOINT`. The worker exports OTLP over HTTP to `OTEL_EXPORTER_OTLP_ENDPOINT`. Both use `OTEL_SERVICE_NAME`, `OTEL_SERVICE_VERSION`, `OTEL_DEPLOYMENT_ENVIRONMENT`, and optional `OTEL_SERVICE_INSTANCE_ID`.

Production must provide both OTLP endpoint variables to a collector reachable from the API and worker containers. The repository intentionally does not deploy the local LGTM containers in `docker-compose.prod.yaml`; run the collector, Loki, Tempo, Mimir, and Grafana as separately managed monitoring infrastructure with persistent storage and retention appropriate for production.

`infra/monitor/alerts/yjs-alert-rules.yaml` contains the initial PromQL rules. Load it through the production Mimir ruler or an equivalent alerting system, then attach the configured notification policy before relying on the alerts.
