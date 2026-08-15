# ClientGateway API version comparison

Only ClientGateway v1 is currently published, so there is no second public contract to compare yet.

| Capability | v1 Beta |
| --- | --- |
| HTTP contract | OpenAPI 3.1 |
| GraphQL contract | Bundled SDL and operation examples |
| Authentication | HttpOnly access/refresh cookies plus CSRF |
| API keys | User-owned keys can be created, listed, and revoked through the authenticated ClientGateway surface |
| Importable client | Postman Collection 2.1 |

When another API version is introduced, this file must compare base paths, authentication, renamed or removed operations, schema changes, error behavior, rate limits, and migration deadlines.
