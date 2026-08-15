# Notezy ClientGateway v1 public contract

This directory contains the versioned contract for the web/client surface exposed by ClientGateway.

ClientGateway uses the existing account/password and HttpOnly access/refresh cookie flow. It includes authentication, user/account operations, resource operations, GraphQL, realtime ticket issuance, notification access, and static assets that are intentionally not part of APIGateway's external integration contract.

- **Canonical contract:** `openapi/openapi.json` (OpenAPI 3.1)
- **Rules:** `rules/`
- **Endpoint catalog:** `reference/endpoints.md`
- **GraphQL SDL:** `graphql/schema.graphql`
- **Runnable examples:** `examples/curl/` and `examples/http/`
- **Postman:** import the collection and environment files in `postman/`
- **Version records:** `versions/dev-log.md` and `versions/comparison.md`

The API key management routes are also part of the authenticated ClientGateway surface:

- `POST /api/development/v1/me/api-keys/create`
- `GET /api/development/v1/me/api-keys/`
- `DELETE /api/development/v1/me/api-keys/{api-key-id}`

An API key secret is returned only once when created. ClientGateway documentation never contains real credentials.

APIGateway's separate external integration contract is located at [`contracts/api-gateway/v1/public/`](../../../api-gateway/v1/public/).
