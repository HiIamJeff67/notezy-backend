# Microservice Architecture Migration

## Purpose

Notezy adopts a staged migration rather than a one-time domain/database split. `Core` remains the synchronous owner of the current user, shelf, block, routine, and station workflows until a later bounded-context decision proves a separate runtime and database are warranted.

## Target ownership

```text
cmd/
  api/
    commands/
contracts/
  graphql/
  api/v1/                  # public client-facing DTOs
  core/v1/                 # Core-owned internal boundary
  durablejob/v1/           # DurableJob-owned internal boundary
  email/v1/                # Email-owned internal boundary
  yjsworker/v1/            # YjsWorker-owned internal boundary
internal/
  exceptions/
  gateway/
    transports/
      api/
        binders/
        controllers/
        cookies/
        interceptors/
        middlewares/
        routes/
      core/
        adapters/
    server/
  shared/
  platform/
  services/
    core/
    durablejob/             # independent runtime; currently shares PostgreSQL
    email/                  # independent runtime and SMTP sender
    yjsworker/              # external worker boundary
docs/
infra/
tests/
```

Directories are created only when they receive an owned implementation. The first migration baseline moves the existing Go application under `internal/`, creates `cmd/api/` for the API executable and its Cobra commands, and moves GraphQL SDL/documents to `contracts/graphql/`. Later issues establish the remaining owners; no empty directory tree is committed merely to mirror this diagram.

## Dependency direction

```text
cmd/* -> internal/*
gateway -> contracts + shared + platform
gateway -X-> service data/repository
services/* -> contracts + shared + platform + own data
services/* -X-> other service source
shared -X-> exceptions + application/framework packages
platform -X-> domain business packages
```

`internal/shared/lib` is portable: it may use the standard library and a necessary third-party library, but never imports a Notezy package. Shared parsers return ordinary Go `error`; application boundaries map those errors to service exceptions.

## Public and internal transport

Client-facing HTTP uses `routes -> binders -> controllers -> Core service
adapter`. Binders own URI/query/header/body binding and transport validation;
controllers invoke an adapter and render the client-safe response. Client-only
cookies, middlewares, and interceptors remain under the same `api` transport.
Neither layer owns transactions, business workflows, or persistence queries.

Gateway internal-API adapters are outbound clients reusable by REST, GraphQL,
and WebSocket. They live in
`internal/gateway/transports/core/adapters/`; client transport code must
not contain Core service client implementation. Core transports are
inbound adapters: they verify the internal delegation credential, map a
versioned request to a Core service call, and map the service result to a
versioned response. Core adapters exist only for Core outbound calls to another
runtime.

Core's inbound Gateway transport is organized as
`internal/services/core/transports/gateway/endpoints/` and
`internal/services/core/transports/gateway/routers/`. Endpoints own the
delegated request/service/response flow; routers only bind HTTP paths and
methods. Delegation and Core-owned authorization middleware belongs in
`internal/services/core/transports/gateway/middlewares/`. Tests live beside
their endpoint, router, or middleware target.

`cmd/api` is the temporary development composition root: it starts Core, waits
until Core has connected its owned dependencies and begun accepting its private
transport, then starts the Gateway listener. Core never starts the public Gin
listener itself. This keeps the two server lifecycles independently movable to
their future commands without reintroducing a Core-to-Gateway source import.

NOT-57 adds `cmd/durablejob` and `cmd/email` as independent composition roots.
DurableJob claims and executes task records from its service-owned data package;
when it needs Core-owned projection business logic, it uses the versioned
DurableJob-to-Core internal HTTP endpoint and delegation credential. Email owns
its SMTP sender and queue and exposes only its Core-facing transport. Both
commands initialize observability and stop their workers/HTTP servers on
context cancellation or SIGTERM.

The WebSocket protocol's persistence payloads and capability values are shared
transport types under `internal/shared/types`; Core services therefore do not
import Gateway WebSocket types. Gateway owns socket admission, tickets, leases,
connection state, and worker forwarding. Core owns the services used for
authorization, durable Yjs state, and block projection. The next runtime split
replaces the in-process bridge with the existing versioned Core client while
preserving these protocol types.

The Gateway extracts and cryptographically parses browser access/refresh tokens
without querying Core data, then applies route permission policy. It sends a
short-lived internal delegation credential containing the calling component as
`actor`, an authenticated user's public UUID as the optional `userSubject`,
request/trace identity, and allowed permissions. For secure calls it also
forwards browser cookies, User-Agent, and sanitized forwarding metadata through
the private client. Core validates both the delegation credential and the
forwarded browser token, then resolves its own internal user identity. Core
inbound routes use `DelegationMiddleware` for component-only operations and
`DelegationAuthenticatedMiddleware` when a user subject is required. Go
`context.Context` is never serialized as an internal transport contract.

## GraphQL

GraphQL follows Scheme A. The Gateway owns GraphQL execution, resolvers, and dataloaders. Resolvers call Core transports through Gateway adapters; they never access Core repositories, GORM schemas, or data packages directly.

GraphQL SDL, fragments, and operation documents are source contracts in `contracts/graphql/`. Generated Go artifacts, generated models, and scalar infrastructure belong in `internal/platform/graphql/` once GraphQL is migrated. Generated files are never edited directly.

## Exceptions and lifecycle

`internal/exceptions` contains only the base exception envelope and origin
classification. Core domain factories stay in
`internal/services/core/exceptions/`; Gateway-only failures use the base envelope
at their use site. They return the base envelope but are never imported across a
Gateway/service boundary. Gateway transport owns client-safe HTTP exception
rendering. Portable shared libraries and parsers never import or return an
application exception.

Kafka, transactional outbox, consumer reliability, and cross-Gateway event fanout are Phase 3 concerns. Yjs updates, awareness, and presence remain ephemeral transport traffic and do not enter the outbox.
