# ClientGateway v1 development log

## Current contract baseline

- Published surface: 168 ClientGateway operations.
- Contract format: OpenAPI 3.1 with bundled GraphQL SDL.
- Authentication: account/password registration or login followed by HttpOnly cookie reuse and CSRF handling.
- Tooling: Postman 2.1 collection/environment, curl functions, and an HTTP client file.
- Stability: Beta namespace with v1 compatibility rules.

Future entries must identify added, changed, deprecated, and removed operations and link to their migration notes. Never record a breaking change only in application release notes.
