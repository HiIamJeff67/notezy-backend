# Architecture and Ownership

The target structure and staged ownership are defined in
[microservice-architecture.md](../codebase-design/microservice-architecture.md). Do
not create empty packages just to satisfy the target tree; move an owner only
when the corresponding migration issue owns its implementation.

## Public HTTP and internal transport

```text
public HTTP: route -> binder -> controller -> Gateway Core adapter -> Core Gateway endpoint -> Core service -> repository/scope
GraphQL: Gateway executor -> resolver/dataloader -> Gateway Core adapter -> Core Gateway endpoint -> Core service
```

| Layer | Responsibility | Must not do |
| --- | --- | --- |
| Gateway route | Register URL, middleware, trace/metric names, and route permission policy | Parse requests, execute business workflow, query SQL |
| Gateway binder | Bind URI/query/header/body and validate the public HTTP contract before invoking its controller | Decide ownership, open domain transactions, query repositories |
| Gateway controller | Call an adapter with a validated request DTO and render a client-safe response | Bind public HTTP input or query repositories directly |
| Gateway Core adapter | Map a versioned contract to an outbound Core request and map the result/error back | Implement Core business rules or access Core data |
| Core Gateway endpoint | Verify delegation credential, map contract request/response, call Core service, and render the internal HTTP response | Re-apply public route semantics or query GORM directly |
| Core service | Validate workflow request, coordinate transaction and application dependencies | Render Gin/HTTP response |
| Repository/scope | Assemble persistence, permission, preload, soft-delete, and locking query | Import transport request/response or return HTTP status |

Gateway binders own public transport parsing and validation. They live in
`internal/gateway/transports/api/binders/`, alongside controllers. Controllers
receive validated request DTOs, invoke a Gateway adapter, and turn the result
into the client response. They do not contain public HTTP parsing, domain rules,
or data-access code.

Gateway cryptographically parses browser credentials and creates a short-lived
delegation credential for an internal request. The credential carries the calling
component, optional authenticated public user subject, trace/request identity, and
route-declared allowed permissions. Core verifies it, validates the forwarded
browser credential, then resolves its own trusted identity before passing
permissions to service/repository options. Browser JWTs and `context.Context` are
not DTO fields or internal request types.

The client-facing Gateway transport lives in
`internal/gateway/transports/api/`. It owns routes, binders, controllers, and
client-only middlewares/interceptors. Reusable Gin cookie handlers live in
`shared/cookies/`.
`controller_func.go` defines the shared controller function type; binder packages
must import it explicitly as `apitransport`.

Gateway-to-microservice internal HTTP is separate at
`internal/gateway/transports/core/`: its `adapters/` package is the
outbound client boundary. If a Core service concern later needs middleware or
interceptors, it belongs under this transport, never under `api/`. Core's
inbound Gateway transport is
organized by responsibility:

```text
internal/core/transports/
  middlewares/              # delegation, authentication, role, and plan checks
  gateway/
    endpoints/              # endpoint interfaces, handlers, and adjacent tests
    routers/                # HTTP route registration and adjacent router tests
```

An endpoint owns parsing the delegated envelope, invoking a Core service, and
serializing the internal response. A router only constructs groups and binds
HTTP method/path to endpoint methods. Core middleware performs delegation and
Core-owned authorization before an endpoint runs. Keep `endpoint.go` only for
endpoint-wide helpers; render a successful operation response in its endpoint
method.

Kafka is also a transport boundary. Runtime-specific Kafka consumers and
producers belong under the transport owned by the peer they communicate with;
they do not belong in a generic `workers/` package:

```text
internal/core/transports/
  durablejob/
    consumers/
    eventbuilders/
    producers/
  yjsworker/
    consumers/
    producers/
  outbox_relay.go

internal/durablejob/transports/
  core/
    consumers/
    producers/
    strategies/
```

For a peer transport, event directions are kept one-to-one at the file
boundary: each Core-to-DurableJob consumer has a matching DurableJob Core
producer, and each Core producer has a matching DurableJob consumer. A file
under `producers/` is a broker producer only when it owns a `Produce()` method
and receives the platform Kafka producer. A file under `eventbuilders/` only
builds a versioned envelope through `Build()`; it never publishes to Kafka.
Core events written inside a database transaction use the Core `OutboxRelay` as
their Kafka publisher. The event builder and the outbox relay together provide
the atomic database-to-Kafka handoff, so an event builder must not call the
broker directly.

The transport owns Kafka envelopes, producer/consumer setup, retry and offset
handling, and calls the local service or engine through constructor-injected
dependencies. Scheduling and execution policy remain in the owning runtime's
service or engine; those packages must not import their transport back.

## Runtime-owned workers

Long-lived background loops belong to the owning runtime's `workers/` package,
not to `services/` or a generic `workers/` package under `shared/platform`:

```text
internal/<runtime>/
  services/   # request-scoped business workflows
  workers/    # runtime-owned long-lived loops and reconciliation
```

A runtime-owned worker must:

- define an `XxxWorkerInterface` before its concrete worker when the composition
  root or tests need a replaceable lifecycle boundary;
- define an `XxxWorker` struct and `NewXxxWorker(...)` constructor for explicit
  dependency injection;
- expose `Start(context.Context) func()` and use the returned function for
  cancellation and graceful shutdown;
- accept context cancellation, avoid reading environment variables, and never
  register HTTP routes or render HTTP responses;
- own scheduling concerns such as tickers, bounded scans, retries, and
  reconciliation triggers while delegating business mutations to injected
  services/repositories;
- be constructed and started only by the runtime's application composition root.

Workers are runtime infrastructure with business-domain awareness. They may
coordinate the owning runtime's data and service dependencies, but they must not
import another runtime's source package or become a hidden replacement for a
service method. A service method must remain directly callable without starting
the worker.

Core services are grouped by business ownership so a future runtime split has a
stable package boundary:

```text
internal/core/services/
  auth/                               # auth and OAuth
  user/                               # user, account, info, settings, billing plans
  shelves/                            # root shelf, sub shelf, item
  blocks/                             # block pack, block, Yjs persistence
  material/
  routines/                           # station, routine, routine tag, RoutineTask
    handlers/                          # RoutineTask execution handlers by aggregate
      block_pack_handler.go
      root_shelf_handler.go
      routine_handler.go
      sub_shelf_handler.go
    matchers/                          # template matching
    parsers/                           # RoutineTask payload decoding and flattening
    resolvers/                         # RoutineTask and block-pattern resolution
  other/                              # badge and theme
  realtime/
```

`routines/routine_task_execution_service.go` owns Core's transaction, permission,
and result-application boundary. Aggregate-specific RoutineTask execution belongs
to the `handlers` package (`RootShelfHandler`, `SubShelfHandler`,
`BlockPackHandler`, and `RoutineHandler`), each with an interface and constructor
for dependency injection. Pattern resolution, payload parsing, and template
matching live in their respective `resolvers`, `parsers`, and `matchers`
packages. Handler constructors receive the base `*gorm.DB` and initialize their
own repositories; operation methods receive the exact `*gorm.DB` session from the
execution service, so a transaction and a normal session use the same path and
never require a service clone such as `withTransactionDB`. Handler methods use an
explicit `Handle...` prefix (for example, `HandleCreateBlockPack`) so they remain
visually distinct from Core service methods.
These packages must not contain a second service orchestration layer or transaction
helper.
Pure assignment execution and template interpolation remain in the DurableJob
runtime; Block remains a projection read model and must not gain RoutineTask
append/update/reset mutation methods in Core.

## Dependency direction

- `cmd/*` may import `internal/*`.
- Gateway client/API and Core-adapter transport code may import contracts, shared,
  and its own code; it must not query Core data or import repositories/
  GORM schemas. RealtimeGateway does not construct Core services, query Core
  data, or synchronously call Core after ticket issuance; it communicates with
  YjsWorker and receives Core lifecycle facts through Kafka.
- A runtime may import contracts, shared, and its own data. A runtime
  must not import another service source package.
- Gateway, DurableJob, Email, and RealtimeGateway are separate Go environments:
  each runtime owns a `go.mod` and `go.sum` beside its `application.go`. The
  root module still composes `cmd/*`, contracts, shared, and Core during
  migration; runtime modules may use a local root-module replacement only for
  explicitly tracked transitional dependencies. Do not add a runtime
  dependency merely to reuse another runtime's source. YjsWorker keeps its
  independent Node/TypeScript package environment.
- `shared` is the root-level cross-runtime utility layer. It may depend on
  contracts and the minimum common application support it genuinely needs;
  portable `shared/lib` never imports a Notezy package.
- `shared/cookies`, `shared/exceptions`, and `shared/tokens` are shared semantic
  boundaries that remain at the root of `shared`; reusable implementation
  utilities belong under `shared/util/` (`editableblock`, `exceptionwriter`,
  and `responsewriter`). `shared/util` may use application-support packages,
  while `shared/lib` remains the stricter dependency-free library layer.
- The generic Kafka envelope is maintained in `contracts/types/event.go`.
  Runtime event domains remain under their owning `contracts/<runtime>/v1/events/`
  package; email request payloads therefore live in
  `contracts/email/v1/events/`.
- `shared/platform` owns infrastructure mechanics, not User/Shelf/Routine
  business rules.
- Cross-runtime calls use a versioned contract and adapter/client. Core adapters
  are outbound only; a Core inbound transport is already the inbound adapter.

## Composition roots

Environment configuration follows the same ownership rule. A runtime composition
root loads each typed owner config once from its `configs/` package and injects it into dependencies;
infrastructure config is colocated with its component at
`shared/platform/<component>/config.go`. Do not read environment variables
from transports, services, workers, clients, or middleware, and do not recreate
`shared/platform/config/`.

Runtime-owned Redis cache ranges and TTLs follow the same boundary. Cache
clients receive a cache-specific config from the runtime composition root; cache
client constructors must not embed Redis database range or user-data expiry
values. The resulting client is reused for its runtime's stores, services, and
transports instead of constructing a new client inside an operation.

Do not introduce an application `modules/` package merely to wrap service
construction. The owning composition root constructs its scope -> repository ->
service dependencies directly, then passes each concrete service to the router or
runtime that uses it. Core's `NewCoreTransportRouter` is the composition root for
its inbound endpoints; the WebSocket runtime constructs only its own Gateway,
Core client, Redis lease store, and YjsWorker manager.

Router construction may instantiate endpoint objects from the services it
receives, but it must not recreate the services or conceal their dependencies in a
module wrapper. Constructor parameters, struct dependency fields, and constructor
assignments use the same order.

Do not create a module, interface, adapter, or helper for a single anticipated
future use. A concrete boundary with a real caller is enough.

## Public operation ordering

One public operation has one ordering across route registration, binder
interface/implementation, controller interface/implementation, Gateway adapter,
Core Gateway endpoint interface/implementation, router registration, and Core service interface/implementation:

1. read: get one, get many/search, aggregate;
2. create: one then many;
3. update: one then many;
4. restore: one then many;
5. soft delete: one then many;
6. hard delete: one then many;
7. permission/sub-resource: get, create, update/upsert, delete;
8. visualization/chart;
9. GraphQL, system-only, or background operation.

Visualization/chart operations form their own family. GraphQL, system-only, and
background operations form another family when present. Do not casually reorder a
legacy file that is outside the operation group being changed.

## File and helper layout

Each controller, service, repository, or adapter file uses this order:

```text
package / imports
interface
concrete struct
constructor
optional auxiliary helpers
public methods in interface order
optional visualization/chart methods
optional GraphQL/system-only methods
```

- Keep one blank line between top-level methods.
- Extract a helper only when two or more methods reuse the same named concept, or
  the inline logic would hide the primary workflow. One-call parsing, mapping,
  validation, temporary type, and wrapper variable stay inline.
- Use `sep30` only when a file has two or more independently navigable method
  families. Auxiliary helpers, visualization/chart methods, and GraphQL/system-only
  methods use `sep30` when they coexist with another family. Do not use it between
  ordinary methods or above a file's only method family.
- Local struct/type declarations require a concrete domain name and repeated use
  within a complex query/result mapping. Do not create `Data`, `Result`, or
  `Params` wrappers for one handoff.

## GraphQL and background runtimes

GraphQL uses Scheme A: Gateway owns the executor, resolvers, and dataloaders.
GraphQL source SDL/fragments/documents, generated Go code, scalars, and generated
models live in `contracts/core/v1/graphql`. Generated files are regenerated from source and
never edited directly. GraphQL business
RequestDto/ResponseDto live in the same
`contracts/core/v1/api/<route-domain>/search.go`
as their owning Core service. Core exposes each GraphQL operation from that
service's endpoint and router; never create a shared GraphQLEndpoint or a central
Core GraphQL router.

DurableJob, Email, and YjsWorker own their own runtime, transport, and service-local
data/types/config. Core and every other runtime may add a runtime-owned
`workers/` package for long-lived background coordination. They must support
`context.Context` cancellation and graceful shutdown. Kafka, outbox, and consumer
reliability are separate Phase 3 concerns; Yjs update, awareness, and presence
stay out of the outbox.

## Minimal pre-change check

1. Locate the Gateway route/controller/adapter or Core Gateway endpoint/Core service
   owning the operation.
2. Locate the service data repository/scope and closest existing test.
3. Confirm the target dependency direction before adding an import.
4. Add a new package or boundary only for a real independently owned lifecycle or
   external contract.
