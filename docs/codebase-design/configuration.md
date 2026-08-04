# Runtime Configuration

Configuration is owned by the runtime or infrastructure component that consumes
it. Environment variables are only read by a typed `LoadConfig` function; the
runtime composition root loads each configuration once and injects it into
clients, workers, transports, and services.

## Ownership

| Owner | Config file | Examples |
| --- | --- | --- |
| PostgreSQL connection | `internal/platform/database/config.go` | `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DOCKER_DB_PORT` |
| Redis connection | `internal/platform/redis/config.go` | `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_INIT_DB` |
| Kafka connection and TLS | `internal/platform/kafka/config.go` | `KAFKA_BROKERS`, `KAFKA_DIAL_TIMEOUT`, `KAFKA_TLS_*`, `KAFKA_SASL_*` |
| OpenTelemetry | `internal/platform/observability/config.go` | `OTEL_SERVICE_*`, `OTEL_EXPORTER_OTLP_GRPC_ENDPOINT` |
| Gateway | `internal/gateway/config/config.go` | `GATEWAY_LISTEN_ADDRESS`, `CORE_BASE_URL`, `REALTIME_GATEWAY_BASE_URL` |
| Core | `internal/services/core/config/` | `CORE_LISTEN_ADDRESS`, `OAUTH_GOOGLE_*`, `STORAGE_KEY_SALT`, `OUTBOX_RELAY_*` |
| DurableJob | `internal/services/durablejob/config/config.go` | `DURABLEJOB_LISTEN_ADDRESS`, `YJS_*_WORKER_URL` |
| Email | `internal/services/email/config/config.go` | `EMAIL_LISTEN_ADDRESS`, `SMTP_*`, `NOTEZY_OFFICIAL_*` |
| RealtimeGateway | `internal/realtimegateway/config/config.go` | `REALTIME_GATEWAY_LISTEN_ADDRESS`, `REALTIME_ENABLED`, `YJS_WORKER_URLS` |

`internal/platform/config/` must not be recreated. A platform component owns
only its infrastructure connection configuration; runtime policy remains with
the runtime that uses it.

## Canonical duration names

Duration values use Go duration strings. Do not introduce numeric unit suffixes
such as `_SECONDS`, `_MILLISECONDS`, or `_HOURS`.

```dotenv
KAFKA_DIAL_TIMEOUT=3s
CORE_CLIENT_TIMEOUT=10s
EMAIL_SERVICE_CLIENT_TIMEOUT=5s
REALTIME_GATEWAY_CLIENT_TIMEOUT=3s
KAFKA_CONSUMER_INITIAL_RETRY_BACKOFF=250ms
KAFKA_CONSUMER_MAXIMUM_RETRY_BACKOFF=5s
OUTBOX_RELAY_POLL_INTERVAL=1s
OUTBOX_RELAY_CLAIM_TIMEOUT=30s
OUTBOX_RELAY_INITIAL_BACKOFF=1s
OUTBOX_RELAY_MAXIMUM_BACKOFF=1m
OUTBOX_RELAY_RETENTION=168h
OUTBOX_RELAY_CLEANUP_INTERVAL=1h
```

All credentials, salts, passwords, client secrets, and SASL credentials are
secrets. They are injected by local development tooling, Compose, or the
production secret manager and must never be logged or committed.
