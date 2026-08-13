# Version and change rules

The URL version (`v1`) is the compatibility boundary for request and response contracts. Repository release tags use Semantic Versioning independently.

- Backward-compatible fields and endpoints may be added within v1.
- Clients must ignore unknown response fields.
- Removing or changing the meaning/type of a field requires a new API version or a documented migration window.
- Deprecated operations remain in the OpenAPI document with a removal date before deletion.
- The development namespace indicates Beta stability; it does not remove the v1 compatibility obligation for documented behavior.
- OpenAPI, Postman, examples, rules, routes, and Go DTO changes must ship together.

Compare generated OpenAPI artifacts between releases to produce the public API change log.
