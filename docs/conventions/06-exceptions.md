# Exception Conventions

## Ownership and envelope

- `internal/exceptions` is a pure application envelope. It may use the standard
  library and `internal/shared`, but must not import Gateway, a microservice,
  platform observability, Gin, GORM, or generated GraphQL code.
- `internal/exceptions` creates the shared envelope with `exceptions.New()`.
  Its optional final `isInternal ...bool` accepts only the first value and
  cannot be changed after construction. It must not own a domain factory or a
  numeric code registry.
- Gateway and each microservice may define a local factory in their own
  `exceptions/` package when two or more callers share the same domain error
  semantics. Examples are `internal/gateway/exceptions/` and
  `internal/services/core/exceptions/`. A component with its own operational
  failure semantics may do the same beneath the component, such as
  `internal/shared/storage/exceptions/`.
- A local factory returns `*exceptions.Exception`, contains no numeric code
  registry, and is never imported across a Gateway/service boundary. Gateway
  never imports an API local factory, and one service never imports another
  service's factory. One-off errors, including generic utility failures, use
  `exceptions.New()` at the call site.
- Set `Reason`, `Domain`, and `Operation` explicitly. They are stable,
  machine-readable application properties; `Message` is the human-readable
  explanation. Use `WithOrigin(err)` for the underlying Go error and
  `WithDetails()` only for in-process diagnostics.
- `Exception` never imports observability and has no `Log()` method. At an
  owning runtime boundary, use `logs.NotezyLogger.JSON(ctx, slog.LevelError,
  exception.String(), exception)` to serialize and record the exception. The
  logger accepts `any`, owns JSON marshaling, and returns its marshal error;
  callers that can safely continue may explicitly discard that logging error.
  `Exception.String()` is diagnostic-only and may contain origin/details; it
  must never be sent to a client.
- Root `ExceptionCode`, `ExceptionPrefix`, generic domain helpers, and domain
  factories do not exist. Domain factories are local to their owning runtime
  or component; callers must not recreate a root compatibility layer.

```go
exception := exceptions.New(
	"NotFound",
	"RootShelf",
	"GetMyRootShelfById",
	"Root shelf was not found",
	http.StatusNotFound,
)

exception := exceptions.New(
	"DatabaseFailure",
	"API",
	"CreateRootShelf",
	"Failed to persist the root shelf",
	http.StatusInternalServerError,
	true,
).WithOrigin(result.Error)
```

## Public safety boundary

- Only `internal/shared/responsewriter.ToPublic()` may convert an exception
  into a browser-facing exception. It records the original internal exception
  before conversion and removes `Origin`, `Details`, numeric compatibility
  codes, and all implementation diagnostics.
- An internal exception becomes the safe, generic `InternalServerError` unless
  the owner explicitly provides an approved fallback with
  `WithPublicFallback()`. A fallback may expose only a safe reason, domain,
  operation, message, HTTP status, and retryability.
- Gateway REST and GraphQL response paths must use the response writer. Do not
  construct an exception JSON body directly and do not pass `Origin.Error()` to
  a response or GraphQL extension.
- Gateway records the original exception through the platform logger before
  calling `ToPublic()`. Internal service transport responses carry the same
  safe `reason`, `domain`, `operation`, `message`, and `retryable` fields; a
  Gateway adapter reconstructs that envelope without domain-specific mapping.
- `ToPublic()` is idempotent: calling it for an already public exception keeps
  its safe fields and still drops diagnostics.

```go
exception := exceptions.New(
	"StorageFailure",
	"API",
	"UploadMaterial",
	"storage provider returned an internal response",
	http.StatusInternalServerError,
	true,
).WithPublicFallback(exceptions.PublicFallback{
	Reason:         "ServiceUnavailable",
	Domain:         "Storage",
	Operation:      "UploadMaterial",
	Message:        "Storage is temporarily unavailable",
	HTTPStatusCode: http.StatusServiceUnavailable,
	Retryable:      true,
})
```

## Context ownership

- `internal/shared/contexts` contains only generic `context.Context` value
  helpers; it never imports exceptions or framework code.
- `internal/gateway/contexts` owns Gin/request-context parsing and the
  route-declared permission set used to construct an outbound delegation.
- `internal/services/core/contexts` owns verified API-service context values,
  including permissions reconstructed from the verified delegation credential.
  Core services must not trust client DTO fields for actor identity or route
  permission policy.

## `exceptions.Cover()`

`exceptions.Cover(existing, fallbacks)` is kept for a single operation with two
or more alternative failure conditions. It returns `existing` unchanged, or the
first satisfied fallback when `existing` is nil. Put an underlying result error
before empty-result and invariant fallbacks so the diagnostic cause is retained.

Use a direct `if` for one clear condition. `Cover()` does not log, roll back, or
write a response; return or finish transaction handling immediately after it
selects an exception.
