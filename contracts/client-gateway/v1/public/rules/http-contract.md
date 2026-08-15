# HTTP contract rules

- Base path: `/api/development/v1` for the current Beta namespace.
- Request and response media type: `application/json`, except GraphQL Playground GET.
- Path resource identifiers are UUID strings unless the operation schema says otherwise.
- Times use RFC 3339 date-time strings.
- Public success envelope: `{ "success": true, "data": ..., "exception": null }`.
- Public failure envelope: `{ "success": false, "data": null, "exception": ... }`.
- `exception.retryable` is the server's explicit retry signal. A client must not infer retryability only from the message.
- Optional `embedded.publicId` identifies the authenticated actor.
- Optional `refreshableTokens.newCSRFToken` replaces the previously stored CSRF value.
- Unknown request fields should not be used for forward compatibility. Only documented properties form the contract.
- Batch requests are not atomic unless the operation description or future version explicitly promises atomicity.
- DELETE may be soft delete or permanent delete; permanent endpoints include `permanently` in their path.

The OpenAPI operation extensions `x-go-request-dto` and `x-go-response-dto` identify the source contracts used to generate each schema.
