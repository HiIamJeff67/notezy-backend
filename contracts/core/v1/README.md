# Core service Contract v1

`Request[D]` and `Response[D]` are the versioned envelope for
Gateway-to-service calls. The HTTP delegation credential is carried only in the
`Authorization: Bearer` header; browser JWTs and Go contexts are never
serialized into this envelope.

The Gateway supplies the route operation, request identity, trace context, and
an optional idempotency key. The delegation credential identifies the calling
component in `actor`; an authenticated user is carried separately as the
optional public UUID `userSubject`. Secure Gateway calls also forward the
browser access/refresh cookies and sanitized request identity headers so Core's
authentication middleware can validate the browser credential.

Core transport uses `DelegationMiddleware` for operations that do not require a
user and `DelegationAuthenticatedMiddleware` for operations that do.

If Core authenticates through the forwarded refresh cookie, it returns rotated
tokens only through private response headers (`X-Core-Auth-Refreshed`,
`X-Core-Set-Access-Token`, and `X-Core-Set-CSRF-Token`). The Gateway consumes
those headers and updates the browser-facing response; Core never writes client
cookies directly.

Internal clients do not retry writes unless an idempotency key is present and
the operation-specific Core transport explicitly supports it. This initial v1
client performs no automatic retries.

`CORE_LISTEN_ADDRESS` configures the private Core listener (default
`127.0.0.1:8080`). Gateway resolves it through `CORE_BASE_URL` (default
`http://127.0.0.1:8080`). A deployment with separate Gateway and Core processes
sets the latter to the Core service's private network address; neither listener
is registered through the browser-facing route tree.
