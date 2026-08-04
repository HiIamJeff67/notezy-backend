# Microservice Architecture Migration

## Purpose

Notezy adopts a staged migration rather than a one-time domain/database split. `Core` remains the synchronous owner of the current user, shelf, block, routine, and station workflows until a later bounded-context decision proves a separate runtime and database are warranted.

## Target ownership

```text
cmd/
  api/
  core/
  realtimegateway/
  durablejob/
  email/
contracts/
  gateway/v1/              # Gateway-owned private request/response envelope
  core/
    v1/
      api/                 # Core-owned HTTP and GraphQL operation DTOs
      events/              # Core-owned Kafka lifecycle event contracts
      graphql/             # Core-owned GraphQL schema and generated artifacts
      types/               # Core-only reusable DTO shapes
  types/
    enums/                 # Canonical cross-runtime enum values
    editable_block.go      # Arborized input tree and raw flattened persistence shape
    *_routine_task_payload.go
  durablejob/v1/           # DurableJob-owned internal boundary
  email/v1/                # Email-owned internal boundary
  realtime-gateway/v1/     # RealtimeGateway-owned internal boundary
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
      realtimegateway/
        adapters/
  realtimegateway/           # standalone WebSocket edge runtime
    data/cache/              # Realtime leases and Realtime-specific rate limits
    transports/
      gateway/               # private API Gateway presence transport
      websocket/             # WebSocket client, connection/channel state, middleware
    workers/                 # YjsWorker connection manager
  platform/
    kafka/                     # Kafka client, readiness, and common telemetry
    redis/                     # Redis client lifecycle + RedisCacheStore registry
  services/
    core/
      data/
        database/             # schemas, repositories, scopes, SQL, seeds, options
        cache/                # Core-owned Redis caches and Lua libraries
        storage/              # Core-owned storage implementations
    durablejob/             # independent runtime; currently shares PostgreSQL
    email/                  # independent runtime and SMTP sender
shared/
  cookies/                    # reusable Gin access/refresh cookie handlers
  editableblock/              # tree-to-row EditableBlock conversion
  constants/ contexts/ responsewriter/ tokens/ types/ validations/
  lib/                       # portable libraries, including BlockNote schemas
yjs-worker/                 # standalone TypeScript runtime; no src/ layer
docs/
infra/
tests/
```

Directories are created only when they receive an owned implementation. The first migration baseline moves the existing Go application under `internal/`, creates `cmd/api/` for the API executable and its Cobra commands, and moves GraphQL SDL/documents to `contracts/core/v1/graphql/`. Later issues establish the remaining owners; no empty directory tree is committed merely to mirror this diagram.

## Dependency direction

```text
cmd/* -> internal/*
gateway -> contracts + shared + platform
gateway -X-> service data/repository
services/* -> contracts + shared + platform + own data
services/* -X-> other service source
shared/lib -X-> all Notezy packages
platform -X-> domain business packages
```

`shared/lib` is portable: it may use the standard library and a necessary third-party library, but never imports a Notezy package. `shared` utility packages such as `responsewriter` may depend on application-facing packages when they provide a deliberately shared runtime utility. Shared parsers return ordinary Go `error`; application boundaries map those errors to service exceptions.

## Shared contract types

`contracts/types/enums` is the canonical owner of cross-runtime enum values.
Core and DurableJob database enum wrappers import those values and add only the
PostgreSQL responsibilities (`Name`, `Scan`, `Value`, validation, and string
conversion); neither runtime redefines a value set.

`contracts/types/editable_block.go` has exactly two EditableBlock shapes:
`ArborizedEditableBlock` is the recursive, validated transport tree, while
`RawFlattenedEditableBlock` is the persistence/projection row with parent and
sibling identifiers plus raw JSON props/content.
`shared/editableblock.FlattenEditableBlock` and
`shared/editableblock.FlattenEditableBlocks` convert the former to the latter. There are
no generic `ArborizedBlock`, `RawArborizedEditableBlock`, or duplicate
`FlattenedEditableBlock` types.

Cross-runtime RoutineTask payloads live directly under `contracts/types/` by
domain. Core parses them into a task request; DurableJob handlers consume the
same contract without importing Core source. BlockNote schemas remain portable
in `shared/lib/blocknote`.

## Public and internal transport

Client-facing HTTP uses `routes -> binders -> controllers -> Core service
adapter`. Binders own URI/query/header/body binding and transport validation;
controllers invoke an adapter and render the client-safe response. Client-only
cookies, middlewares, and interceptors remain under the same `api` transport.
Neither layer owns transactions, business workflows, or persistence queries.

Gateway internal-API adapters are outbound clients reusable by REST and
GraphQL. They live in
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

`cmd/api`, `cmd/core`, and `cmd/realtimegateway` are independent composition
roots. API Gateway accepts HTTP and GraphQL browser traffic only; Core owns
PostgreSQL-backed operations and its private listener; RealtimeGateway is the
separate WebSocket edge runtime, directly addressed by the Nginx WebSocket
upstream. RealtimeGateway owns connection admission, ticket verification, Redis
leases, Realtime-specific rate limits, and long-lived Go-to-YjsWorker
connections. Its root `application.go` is the composition root: each `Start()`
creates a new WebSocket client with its own connection state; no
global singleton application instance is retained. It never constructs Core
repositories or services.

NOT-57 adds `cmd/durablejob` and `cmd/email` as independent composition roots.
DurableJob claims and executes task records from its service-owned data package;
when it needs Core-owned projection business logic, it uses the versioned
DurableJob-to-Core internal HTTP endpoint and delegation credential. Email owns
its SMTP sender and queue and exposes only its Core-facing transport. Both
commands initialize observability and stop their workers/HTTP servers on
context cancellation or SIGTERM.

RealtimeGateway owns socket admission, tickets, leases, connection state, and
worker forwarding. Core owns authorization, durable Yjs state, and block
projection, but RealtimeGateway does not synchronously call Core for those
operations. YjsWorker sends its persistence and projection commands to Core
through Kafka; Core publishes lifecycle facts through its transactional outbox.

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

The versioned Core `Response[D]` envelope carries only operation data. A
short-lived signed BlockPack channel ticket carries the room-admission policy
snapshot required by the WebSocket runtime: policy version, admission enforcement
strategy, and maximum subscribers. RealtimeGateway verifies that ticket and uses
the attested strategy and maximum-subscriber claim directly for atomic Redis lease
admission; it does not persist a Core-owned room-policy cache or expose internal
instructions through the public API.

## GraphQL

GraphQL follows Scheme A. The Gateway owns GraphQL execution, resolvers, and dataloaders. Resolvers call Core transports through Gateway adapters; they never access Core repositories, GORM schemas, or data packages directly.

GraphQL SDL, fragments, operation documents, generated Go artifacts, generated models, and scalar infrastructure belong in `contracts/core/v1/graphql/`. Generated files are never edited directly.

## Data cache ownership

Redis is divided into three explicit responsibilities: platform client lifecycle,
runtime-owned cache registration, and domain-specific cache operations.

```mermaid
flowchart TB
  subgraph Platform["internal/platform/redis"]
    Manager["DefaultClientManager<br/>map[int]*redis.Client"]
    Registry["RedisCacheStores<br/>map[int]RedisCacheStore"]
  end

  subgraph Core["Core runtime"]
    CoreStart["core.Start()"]
    UserRegister["userdata.Register()"]
    UserStore["UserDataCacheStore<br/>DB 0–3"]
    UserClient["UserDataCacheClient"]
    RealtimeRegisterCore["realtimelease.Register()"]
  end

  subgraph Gateway["Gateway runtime"]
    GatewayStart["gateway.Start()"]
    RateRegister["ratelimitrecord.Register()"]
    RateStore["RateLimitRecordCacheStore<br/>DB 4–7"]
    RateClient["RateLimitRecordCacheClient"]
  end

  subgraph RealtimeGateway["RealtimeGateway runtime"]
    RealtimeGatewayStart["realtimegateway.Start()"]
    RealtimeRegisterGateway["realtimelease.Register()"]
    LeaseStore["RealtimeLeaseCacheStore<br/>realtime DB"]
    LeaseClient["RealtimeLeaseStore"]
  end

  CoreStart --> UserRegister --> Manager
  UserRegister --> UserStore --> Registry
  UserClient --> Registry
  RealtimeRegisterCore --> Manager

  GatewayStart --> RateRegister --> Manager
  RateRegister --> RateStore --> Registry
  RateClient --> Registry

  RealtimeGatewayStart --> RealtimeRegisterGateway --> Manager
  RealtimeRegisterGateway --> LeaseStore --> Registry
  LeaseClient --> Registry
```

### Platform lifecycle

`internal/platform/redis` owns no cache domain or business policy. Its
`ClientManager` is the sole owner of Redis connection creation, lookup, and
shutdown for the current process. It maintains `map[int]*redis.Client`, keyed
by Redis database number.

`RedisCacheStore` is the initialization boundary for a concrete database store:

```go
type RedisCacheStore interface {
    DatabaseNumber() int
    Initialize(ctx context.Context) error
}
```

`RegisterCacheStores()` first calls each store's `Initialize()` and only
registers successful stores in `RedisCacheStores`. Consequently, a failed Lua
library load never leaves an apparently usable cache store in the registry.

Each executable has its own process-local `DefaultClientManager` and
`RedisCacheStores` registry. Core, Gateway, and WebSocket therefore have their
own Redis TCP connections even when they target the same Redis server and DB
numbers.

### Runtime ownership

| Runtime | Registration | Redis ownership |
|---|---|---|
| Core | `userdata.Register()` | User data and quota cache, DB 0–3 |
| Gateway API | `ratelimitrecord.Register()` | API rate-limit records, DB 4–7 |
| RealtimeGateway | `realtimelease.Register()` + `ratelimitrecord.Register()` | Connection/channel leases, participant presence, Realtime rate-limit records, lifecycle Pub/Sub fan-out |
| DurableJob / Email | None currently | No registered Redis cache domain |

`UserDataCacheStore` and `RateLimitRecordCacheStore` each own a single
`*redis.Client` and their cache library loading. Their matching `*CacheClient`
types own key formatting, routing, and cache operations; they retrieve a store
from the platform registry rather than keeping a Redis client map themselves.

Core quota functions are one-function-per-file under
`internal/services/core/data/cache/userdata/libraries/`. The UserData store
embeds and joins them into one `user_quota_library`, then performs a single
`FUNCTION LOAD REPLACE` during `Initialize()`.

### Registration and operation flow

```text
core.Start()
  -> userdata.Register(ctx, DefaultClientManager)
    -> ConnectAll(DB 0–3)
    -> NewUserDataCacheStore(DB, client)
    -> RegisterCacheStores(stores...)
      -> store.Initialize()
        -> FUNCTION LOAD REPLACE user_quota_library
      -> RedisCacheStores[DB] = store

Auth/User service
  -> UserDataCacheClient.Get/Set/Update(...)
    -> hash identifier -> database number
    -> GetRedisCacheStore(databaseNumber)
    -> UserDataCacheStore.redisClient
```

`internal/realtimegateway/data/cache/realtimelease/` owns the RealtimeGateway
runtime's user
connection and BlockPack subscriber lease lifecycle, active lease inspection,
participant presence, and presence PubSub fanout. A Core-issued BlockPack
channel ticket carries the signed policy version, reject-new-subscriber strategy,
and maximum-subscriber policy snapshot; RealtimeGateway uses those verified claims
directly during atomic admission.
Core authorizes channel admission when it issues a signed ticket, but it does
not participate in RealtimeGateway subscribe or write boundaries and does not
read RealtimeGateway Redis or the live subscriber count during ownership or
plan workflows. RealtimeGateway applies the ticket's verified channel policy;
its Lua scripts are private to that cache domain.

RealtimeGateway owns the lease cache and produces participant snapshots and
deltas. Its lifecycle Kafka consumer publishes received revocations to its local
Redis Pub/Sub channel so every RealtimeGateway instance detaches only its own
matching connections and releases their leases. The API Gateway REST participant
controller requests that snapshot through the versioned private
`contracts/realtime-gateway/v1` transport, then calls Core only to validate the
requester's BlockPack permission and enrich a bounded, de-duplicated public-ID
collection after a RootShelf membership filter. Core's lifecycle path is outbox
→ Kafka → RealtimeGateway; it never reads or writes RealtimeGateway-owned
lease/cache state.

`internal/services/core/data/storage/` owns Core's storage implementation;
Gateway does not access it directly.

## Lifecycle event contracts

Cross-runtime lifecycle events are defined in
[`contracts/core/v1/events`](../../contracts/core/v1/events/) and documented in
[Kafka Event Contracts](../system-design/kafka-event-contracts.md). Core owns
the facts and policy values in those messages. RealtimeGateway consumes them to
execute detach or future admission behavior, but never derives business rules.
The event envelope is independent of the Kafka client, outbox schema, and
consumer implementation so NOT-34, NOT-35, and NOT-33 can evolve those owners
without changing the semantic boundary.

## Exceptions and lifecycle

`internal/exceptions` contains only the base exception envelope and origin
classification. Core domain factories stay in
`internal/services/core/exceptions/`; Gateway-only failures use the base envelope
at their use site. They return the base envelope but are never imported across a
Gateway/service boundary. Gateway transport owns client-safe HTTP exception
rendering. Portable shared libraries and parsers never import or return an
application exception.

Kafka, transactional outbox, consumer reliability, and cross-Gateway event
fanout are Phase 3 concerns. Core writes a service-owned outbox record in the
same transaction as a lifecycle mutation; WebSocket consumes the resulting
event and only executes the already decided detach/policy action. Yjs updates,
awareness, presence, and Redis leases remain ephemeral transport traffic and do
not enter the outbox.

The platform Kafka client owns broker configuration, TLS/SASL setup, broker
readiness checks, and common telemetry names. Core and RealtimeGateway may run
in degraded mode when Kafka is unavailable, but their readiness endpoints stay
unready until a broker ping succeeds. Local provisioning and operational steps
are documented in [Kafka Local Development](../runbooks/kafka-local-development.md).

Core's PostgreSQL-to-Kafka handoff is a Core-owned transactional outbox. Domain
transactions persist `OutboxEventTable` rows before committing; a Core worker
claims rows with `FOR UPDATE SKIP LOCKED`, publishes the versioned event
envelope, and marks acknowledgements. This deliberately provides at-least-once
delivery, with consumer idempotency owned separately. See the
[transactional outbox design](../system-design/transactional-outbox.md).
