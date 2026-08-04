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
internal/services/core/transports/
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

## Dependency direction

- `cmd/*` may import `internal/*`.
- Gateway client/API and Core-adapter transport code may import contracts, shared,
  platform, and its own code; it must not query Core data or import repositories/
  GORM schemas. RealtimeGateway does not construct Core services, query Core
  data, or synchronously call Core after ticket issuance; it communicates with
  YjsWorker and receives Core lifecycle facts through Kafka.
- A service may import contracts, shared, platform, and its own data. A service
  must not import another service source package.
- `shared` is the root-level cross-runtime utility layer. It may depend on
  contracts and the minimum common application support it genuinely needs;
  portable `shared/lib` never imports a Notezy package.
- `internal/platform` owns infrastructure mechanics, not User/Shelf/Routine
  business rules.
- Cross-runtime calls use a versioned contract and adapter/client. Core adapters
  are outbound only; a Core inbound transport is already the inbound adapter.

## Composition roots

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
data/types/config. They must support `context.Context` cancellation and graceful
shutdown. Kafka, outbox, and consumer reliability are separate Phase 3 concerns;
Yjs update, awareness, and presence stay out of the outbox.

## Minimal pre-change check

1. Locate the Gateway route/controller/adapter or Core Gateway endpoint/Core service
   owning the operation.
2. Locate the service data repository/scope and closest existing test.
3. Confirm the target dependency direction before adding an import.
4. Add a new package or boundary only for a real independently owned lifecycle or
   external contract.
