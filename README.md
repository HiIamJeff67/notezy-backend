<a><img src="global/images/logo/NotezyDocumentationHeaderImage.png" alt="notezy" /></a>

## Contents

- [Project Overview](#project-overview)
- [Core Features](#core-features)
- [Runtime Mechanisms](#runtime-mechanisms)
- [Layered Architecture](#layered-architecture)
- [Build and Run](#build-and-run)
- [Infrastructure and Operations](#infrastructure-and-operations)
- [VS Code Workspace Configuration](#vs-code-workspace-configuration)
- [Key Highlights and Engineering Notes](#key-highlights-and-engineering-notes)
- [Repository Structure](#repository-structure)
- [Licensing and Third-Party Notices](#licensing-and-third-party-notices)

### Project Overview

- Backend stack: Go `1.26`, Gin, GORM, Cobra CLI.
- API style: REST + GraphQL.
- Domain focus: auth/user/account/profile/settings + shelf/material/block + routine/station workflows.
- Storage and state: PostgreSQL + Redis.
- Observability: OpenTelemetry + LGTM stack (Loki/Tempo/Mimir/Grafana).

### Core Features

- Token-based authentication pipeline with cookies and CSRF verification.
- Modular service boundaries with clear binder/controller/service/repository separation.
- Permission-aware data access in repositories via scopes and access-control checks.
- Database lifecycle support: enum migration, table migration, triggers, constraints, seed data.
- Rate limiting for both anonymous and authenticated requests.
- Async email workflow with retry and priority scheduling.
- GraphQL resolver layer with generated schema bindings and dataloaders.

### Runtime Mechanisms

- **Email Worker Manager**
  - `app/emails/manager.go`
  - Priority queue + worker pool + retry/backoff + queue monitor ticker.
  - Task types include welcome/validation/security notifications.
- **Timeout Guard**
  - `app/middlewares/timeout_middleware.go`
  - Buffered response writer + timeout context + panic capture for safe timeout responses.
- **Authentication and Token Refresh**
  - `app/middlewares/auth_middleware.go`
  - Access token validation first; fallback to refresh token path.
  - On refresh success, issues new access/csrf token and updates cache/context.
- **CSRF Validation**
  - `app/middlewares/csrf_middleware.go`
  - Validates `X-CSRF-Token` and rotates token when near expiry.
- **Rate Limiting**
  - `app/middlewares/unauthorized_rate_limit_middleware.go`
  - `app/middlewares/authorized_rate_limit_middleware.go`
  - Uses hybrid rate limiting strategy with Redis-backed coordination.
- **Response Interceptors**
  - `ShareableResponseWriterInterceptor` for post-controller response rewriting.
  - `RefreshTokenInterceptor` embeds refreshed tokens in response body/cookie.
  - `EmbeddedInterceptor` appends additional authorized fields (e.g., `publicId`).
- **Request Safety and Network Controls**
  - CORS handling, origin/referer whitelist, `X-Forwarded-For` sanitization, max body size control.
- **Validation and Input Safety**
  - Custom validators for account/password strength, URL scheme allow/deny list, timezone, block content schema.

### Layered Architecture

#### Request Flow

1. `routes` assemble middleware chains and binders.
2. `binders` map HTTP request/context into request DTOs.
3. `controllers` orchestrate service invocation and response format.
4. `services` execute business workflows, transaction boundaries, cache/email/token interactions.
5. `repositories` encapsulate persistence, scopes, permission checks, soft-delete policies, SQL ops.
6. `models/schemas` define DB tables/enums/triggers/constraints/seeds.

#### Layer Breakdown

- **Models Layer**
  - `app/models/schemas`: table schemas.
  - `app/models/schemas/enums`: enum definitions and migration mappings.
  - `app/models/schemas/triggers`: SQL triggers (cascading, projection, accounting, maintenance).
  - `app/models/schemas/constraints`: SQL constraints/indexes.
  - `app/models/seeds`: seed SQL sets for billing and plan limitations.
  - `app/models/inputs`: create/update input contracts.
- **Repository Layer**
  - `app/models/repositories`: persistence APIs.
  - `app/models/scopes`: permission, preload, soft-delete filtering logic.
  - `app/models/sqls`: raw SQL units for targeted operations.
  - `app/options/repository_option.go`: DB/session/transaction behavior control.
- **Service Layer**
  - `app/services`: business orchestration and workflow composition.
  - Integrates DTOs, cache operations, token ops, and async email dispatch.
- **Controller Layer**
  - `app/controllers`: HTTP response boundaries (`success/data/exception`).
- **Binder Layer**
  - `app/binders`: JSON/query/context parsing and DTO assembly.
- **Routes Layer**
  - `app/routes/developmentroutes`: main API routes (`/api/development/v1`).
  - `app/routes/testroutes`: test-only route registration.
- **Commands and CLI**
  - `app/commands`: migration, seed, enum/db inspect, truncation commands.
- **Exception System**
  - `app/exceptions`: centralized error taxonomy + safe response + metric integration.

### Build and Run

#### Prerequisites

- Go `1.26.x`
- Docker and Docker Compose
- `.env` configured for DB/Redis/JWT/CSRF/SMTP/OAuth values

#### Local Run (Host)

1. Configure `.env`.
2. Run `go run main.go`.

#### Development Run (Compose)

1. `docker compose up -d --build`
2. DB migration: `make migrate-hotreload-db`
3. Seed defaults: `make seed-hotreload-db`

#### Production-Like Run

1. `docker compose -f docker-compose.prod.yaml up -d --build`

#### GraphQL Generation

- Generate: `make gql-generate`
- Clean + regenerate: `make gql-regenerate`

### Infrastructure and Operations

- **Containerization**
  - `docker-compose.yaml`: dev topology with API + DB + Redis + Nginx + LGTM.
  - `docker-compose.prod.yaml`: leaner production-like topology.
  - `infra/docker/Dockerfile.dev`: hot-reload image (`air`).
  - `infra/docker/Dockerfile.prod`: multi-stage production build.
- **Nginx Reverse Proxy**
  - `infra/nginx/default.dev.conf`: HTTP reverse proxy to `notezy-api:7777`.
  - `infra/nginx/default.prod.conf`: HTTPS termination + redirect + proxy headers.
- **GraphQL Artifacts**
  - Schemas: `shared/graphql/schemas/**/*.graphql`.
  - Generator config: `infra/graphql/gqlgen.yaml`.
  - Generated outputs: `app/graphql/generated/*_generated.go`, `app/graphql/models/models_gen.go`.
- **Observability (LGTM + OTEL)**
  - Collector: `infra/monitor/otel-collector-config.dev.yaml`.
  - Loki/Tempo/Mimir configs: `infra/monitor/*.yaml`.
  - Grafana datasource provisioning: `infra/monitor/grafana/datasources.dev.yaml`.
  - App-side tracing/metrics initialized in `app/app.go` and route middleware.
- **Scope Exclusion (PayPal)**
  - `infra/paypal/*` and related docs are intentionally excluded from this main infra analysis section.

### VS Code Workspace Configuration

- `/.vscode/settings.json`
  - Format on save enabled.
  - Go formatter pinned to `golang.go`.
  - OTEL collector YAML schema mapping enabled.
  - Relative line numbers enabled.
- `/.vscode/*.code-snippets`
  - Team snippets for controller/service/schema/exception/comment templates.
- `/.vscode/cspell.json`
  - Domain-specific dictionary tuning for project vocabulary.

### Key Highlights and Engineering Notes

- High structural consistency across business modules:
  - Module, controller, and binder boundaries are aligned for predictable maintenance.
- Persistence and policy boundaries are explicit:
  - Permission checks are scope-driven, not ad hoc in controllers.
- Response lifecycle is extensible:
  - Interceptor strategy supports token refresh embedding and response augmentation.
- Redis partitioning strategy exists for differentiated purposes:
  - user-data cache vs rate-limit records mapped to different Redis DB ranges.
- Current testing profile:
  - e2e tests are focused on auth flows.
  - unit tests target utility-level behavior and helpers.

### Repository Structure

```text
go-start-monolithic-kit/
├── app/
│   ├── adapters/
│   ├── binders/
│   ├── caches/
│   ├── commands/
│   ├── configs/
│   ├── contexts/
│   ├── controllers/
│   ├── cookies/
│   ├── dtos/
│   ├── emails/
│   ├── exceptions/
│   ├── graphql/
│   │   ├── dataloaders/
│   │   ├── generated/
│   │   ├── models/
│   │   ├── resolvers/
│   │   └── scalars/
│   ├── interceptors/
│   ├── middlewares/
│   ├── models/
│   │   ├── inputs/
│   │   ├── repositories/
│   │   ├── scopes/
│   │   ├── schemas/
│   │   │   ├── constraints/
│   │   │   ├── enums/
│   │   │   └── triggers/
│   │   ├── seeds/
│   │   └── sqls/
│   ├── modules/
│   ├── monitor/
│   │   ├── logs/
│   │   ├── metrics/
│   │   └── traces/
│   ├── options/
│   ├── routes/
│   │   ├── developmentroutes/
│   │   └── testroutes/
│   ├── services/
│   ├── storages/
│   ├── tokens/
│   ├── util/
│   └── validation/
├── infra/
│   ├── docker/
│   ├── graphql/
│   ├── monitor/
│   ├── nginx/
│   └── paypal/ (excluded scope)
├── shared/
│   ├── constants/
│   ├── graphql/
│   ├── lib/
│   └── types/
├── test/
│   ├── e2e/
│   └── unit/
├── docs/
├── LICENSES/
├── docker-compose.prod.yaml
├── docker-compose.yaml
├── go.mod
├── go.sum
├── LICENSE(tw).md
├── LICENSE.md
├── main.go
├── Makefile
└── README.md
```

### Licensing and Third-Party Notices

- **Project License (Proprietary/EULA)**
  - `LICENSE.md` (English)
  - `LICENSE(tw).md` (Traditional Chinese)
- **Third-Party License Bundle**
  - Summary index: `LICENSES/THIRD_PARTY_NOTICES.txt`
  - Package-level texts: `LICENSES/<domain>/<package>/...`
- **Compliance Position (Current State)**
  - The repository follows a dual model:
    - proprietary EULA for project code
    - bundled third-party notices/licenses for dependencies
  - Even if the repository becomes private, retaining third-party notices is still recommended for internal compliance traceability.
