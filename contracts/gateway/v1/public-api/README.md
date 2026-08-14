# Deprecated Gateway v1 artifact

This legacy output is retained temporarily for migration compatibility. It is no longer the public contract and must not be linked from external documentation. Use `contracts/api-gateway/v1/public/` for the APIGateway integration API and `contracts/client-gateway/v1/public/` for the ClientGateway user/client boundary.

This directory contains the machine-readable and human-readable contract for all 168 versioned routes currently emitted by Gateway v1.

- **Canonical contract:** `openapi/openapi.json` (OpenAPI 3.1)
- **Rules:** `rules/`
- **Endpoint catalog:** `reference/endpoints.md`
- **GraphQL SDL:** `graphql/schema.graphql`
- **Runnable examples:** `examples/curl/all-endpoints.sh` and `examples/http/all-endpoints.http`
- **Cookie session example:** `examples/curl/authenticated-session.sh`
- **Postman:** import both JSON files in `postman/`
- **Version records:** `versions/dev-log.md` and `versions/comparison.md`

The generated artifacts are refreshed from routes and Go DTOs with:

```bash
make -C contracts public-api-gen
```

No real account, password, cookie, CSRF token, or API secret belongs in this directory.
