# Runtime Configuration

Configuration is owned by the runtime or infrastructure component that consumes
it. Environment variables are only read by a typed `LoadConfig` function; the
runtime composition root loads each configuration once and injects it into
clients, workers, transports, and services.

## Ownership

| Owner | Config file | Examples |
| --- | --- | --- |
| PostgreSQL connection | `shared/platform/database/config.go` | `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DOCKER_DB_PORT` |
| Redis connection | `shared/platform/redis/config.go` | `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_INIT_DB` |
| Kafka connection and TLS | `shared/platform/kafka/config.go` | `KAFKA_BROKERS`, `KAFKA_DIAL_TIMEOUT`, `KAFKA_TLS_*`, `KAFKA_SASL_*` |
| OpenTelemetry SDK | `shared/platform/observability/config.go` | `OTEL_SERVICE_*`, `OTEL_EXPORTER_OTLP_GRPC_ENDPOINT` |
| Gateway | `internal/gateway/configs/` | `GATEWAY_LISTEN_ADDRESS`, `CORE_BASE_URL` |
| Core | `internal/core/configs/` | `CORE_LISTEN_ADDRESS`, `OAUTH_GOOGLE_*`, `STORAGE_KEY_SALT`, `OUTBOX_RELAY_*`, user-data cache TTL, Yjs document initialization endpoint/timeout |
| DurableJob | `internal/durablejob/configs/` | `DURABLEJOB_LISTEN_ADDRESS`, runtime Kafka and maintenance strategy settings |
| Email | `internal/email/configs/` | `EMAIL_LISTEN_ADDRESS`, `SMTP_*`, `NOTEZY_OFFICIAL_*`, `KAFKA_*` consumer settings |
| RealtimeGateway | `internal/realtimegateway/configs/` | `REALTIME_GATEWAY_LISTEN_ADDRESS`, `REALTIME_ENABLED`, `YJS_WORKER_URLS` |

`shared/platform/config/` must not be recreated. A platform component owns
only its infrastructure connection configuration; runtime policy remains with
the runtime that uses it.

DurableJob owns the Yjs maintenance strategy policy. Its composition root loads
these values through `internal/durablejob/configs.LoadConfig` and injects one
immutable `YjsMaintenanceStrategyConfig` into the strategy:

```dotenv
DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_PENDING_HINTS=1000
DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_BATCH=32
DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_DISPATCH_WORKERS=8
DURABLEJOB_YJS_MAINTENANCE_MAXIMUM_REQUEST_ATTEMPTS=3
```

## Canonical duration names

Duration values use Go duration strings. Do not introduce numeric unit suffixes
such as `_SECONDS`, `_MILLISECONDS`, or `_HOURS`.

```dotenv
KAFKA_DIAL_TIMEOUT=3s
CORE_CLIENT_TIMEOUT=10s
CORE_USER_DATA_CACHE_EXPIRES_IN=1h
YJS_DOCUMENT_INITIALIZATION_WORKER_TIMEOUT=30s
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
CI/CD credential store or deployment host secure storage and must never be
logged or committed.

## Environment secret storage

The project uses SOPS with age for encrypted environment files. The encrypted
`.env.enc` for development and `.env.<environment>.enc` for other environments,
along with `.sops.yaml`, plaintext environment files, and age identities, are
deployment/local artifacts and are ignored by Git. Transfer encrypted files
only through an approved private channel; never commit them to GitHub. Each
developer, CI system, Jenkins agent, and production host owns a separate age
identity.

New members generate their own identity and send only the public recipient to
the maintainer. The maintainer uses `make env-updatekeys` to add or remove
recipients and `make env-rotate` after a removal or compromise. CI/CD decrypts
only at runtime using its credential-store identity; deployment scripts remove
temporary plaintext files after Compose exits. See
`docs/runbooks/environment-secrets.md` for the complete onboarding and
rotation workflow.

Redis topology is runtime-owned. Each runtime composition root creates an
immutable `shared/platform/redis.ClientSet` and injects it into its cache stores
and clients. Cache clients may select a private shard by hashing a key, but no
server number or cross-runtime Redis registry is part of their configuration.
Core currently reads only `CORE_USER_DATA_CACHE_EXPIRES_IN`; Gateway and
RealtimeGateway keep their rate-limit policies in
`internal/<runtime>/configs/rate_limit.go`.
