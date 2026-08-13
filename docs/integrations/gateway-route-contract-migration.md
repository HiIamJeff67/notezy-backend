# Gateway v1 route contract migration

This is a frontend integration note for the Gateway v1 route contract change. It documents URL compatibility requirements only; no frontend source code is changed by this migration.

## Scope and rollout

The change applies to the development Gateway v1 routes emitted by `internal/gateway/transports/api/routes/developmentroutes`. The old camelCase URLs are removed; clients must deploy the new URLs as one coordinated change. The request DTO and JSON/query field names remain unchanged unless the OpenAPI document says otherwise.

Canonical generated artifacts:

- [`openapi.json`](../../contracts/gateway/v1/public-api/openapi/openapi.json)
- [`endpoints.md`](../../contracts/gateway/v1/public-api/reference/endpoints.md)
- [`Notezy-Gateway-v1.postman_collection.json`](../../contracts/gateway/v1/public-api/postman/Notezy-Gateway-v1.postman_collection.json)
- [`RealtimeGateway OpenAPI`](../../contracts/realtime-gateway/v1/public-api/openapi/openapi.json)

The base URL is `/api/development/v1` in local development. Use the configured production base URL in deployed environments.

## URL naming rules

- Resource groups and action segments use kebab-case: `block-packs`, `send-auth-code`, `register-via-google`.
- Path placeholders use kebab-case: `:block-pack-id`, `:routine-task-id`, `:user-public-id`.
- Query-string and JSON property names continue to follow their DTO contract. For example, a request body may still contain `blockPackId`, while the URL is `/block-packs/{block-pack-id}`.
- Folder names in generated Postman collections and OpenAPI tags use the same API names: `auth`, `users`, `user-info`, `block-packs`, `routine-task-records`, `graphql`, and `static`.
- OpenAPI `operationId` values remain the internal stable operation identifiers; the URL path and Postman request name are the kebab-case values clients should present to users.

## Authentication route changes

| Before | After |
| --- | --- |
| `/auth/registerViaGoogle` | `/auth/register-via-google` |
| `/auth/loginViaGoogle` | `/auth/login-via-google` |
| `/auth/sendAuthCode` | `/auth/send-auth-code` |
| `/auth/validateEmail` | `/auth/validate-email` |
| `/auth/resetEmail` | `/auth/reset-email` |
| `/auth/forgetPassword` | `/auth/forget-password` |
| `/auth/resetMe` | `/auth/reset-me` |
| `/auth/deleteMe` | `/auth/delete-me` |

The ordinary `/auth/register` and `/auth/login` paths are unchanged. `delete-me` remains a destructive operation and must be protected by the same confirmation UX and permission checks as before.

## Path-parameter changes

Replace camelCase placeholders wherever they occur:

| Before | After |
| --- | --- |
| `:blockPackId` | `:block-pack-id` |
| `:blockId` | `:block-id` |
| `:itemId` | `:item-id` |
| `:materialId` | `:material-id` |
| `:parentSubShelfId` | `:parent-sub-shelf-id` |
| `:prevSubShelfId` | `:prev-sub-shelf-id` |
| `:rootShelfId` | `:root-shelf-id` |
| `:routineId` | `:routine-id` |
| `:routineTagId` | `:routine-tag-id` |
| `:routineTaskId` | `:routine-task-id` |
| `:stationId` | `:station-id` |
| `:subShelfId` | `:sub-shelf-id` |
| `:userPublicId` | `:user-public-id` |

Examples:

```text
GET    /block-packs/:block-pack-id
PATCH  /routines/:routine-id
GET    /routine-task-records/routine-task/:routine-task-id
GET    /static/global-images/:id
```

RealtimeGateway follows the same rule. Its presence endpoint is now:

```text
GET /realtime/development/v1/block-pack/:block-pack-id/participants
```

The exact operation list and whether an identifier is a path, query, or body field is defined by OpenAPI. Do not infer a path parameter from a similarly named DTO property.

## Success status change

Every successful operation whose generated operation ID begins with `create` now returns `201 Created`. All other successful Gateway operations continue to return `200 OK`.

The response envelope is unchanged:

```json
{
  "success": true,
  "data": {},
  "exception": null
}
```

Frontend clients should accept `201` as a successful response for create mutations, preserve the response body handling used for `200`, and avoid treating `201` as an error. Error envelopes and status handling remain defined by the generated OpenAPI contract.

## Postman and generated examples

Import the collection and the example environment from `contracts/gateway/v1/public-api/postman/`. Select **Notezy v1 example (no credentials)** as the active environment. Fill `account`, `email`, and `password`, then run `register` or `login` first. Postman stores `accessToken` and `refreshToken` in its cookie jar; the generated test scripts keep the latest CSRF token in `csrfToken`.

The collection uses camelCase environment variable keys such as `blockPackId` only as client-side variable names. Those variables are substituted into kebab-case URL placeholders, so this is expected:

```text
{{blockPackId}}  ->  /block-packs/00000000-0000-4000-8000-000000000001
```

Requests marked `[DESTRUCTIVE]` are intentionally labeled in the generated collection when their HTTP method is `DELETE` or their operation is a reset action. The marker does not change the route, payload, or authorization behavior; it is a review guard for manual execution.

## Frontend migration checklist

- Update route constants and URL builders to kebab-case.
- Update path-parameter template keys in typed route helpers.
- Keep DTO JSON and query keys unchanged unless the OpenAPI schema changes them.
- Treat both `200` and `201` as successful responses, with `201` expected for create operations.
- Refresh generated API types from OpenAPI rather than copying the old route table.
- Update request mocks, browser tests, analytics route allow-lists, and cache keys that contain the old URL spelling.
- Verify authentication and CSRF cookie handling against the generated login/register examples.
- Do not invoke `[DESTRUCTIVE]` routes from automated smoke tests without an isolated test account and explicit cleanup.

Regenerate the contract after route changes with:

```bash
make -C contracts public-api-gen
```
