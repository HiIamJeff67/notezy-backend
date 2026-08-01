# GraphQL source contracts

This directory owns GraphQL source shared with client code generation:

- `schemas/` contains SDL, including scalar and enum definitions.
- `fragments/` and `queries/` contain reusable client operation source.
- `operations.go` contains private Gateway-to-Core operation names. GraphQL
  RequestDto/ResponseDto are owned by the matching
  `contracts/api/v1/<route-domain>/search.go` package.

The Go server generated artifacts, generated models, and scalar implementation
belong to `internal/platform/graphql/`. Gateway GraphQL execution, resolvers, and
dataloaders belong to `internal/gateway/transports/api/graphql/`; they access Core
only through `internal/gateway/transports/core/adapters/CoreClient`. Each Core
domain endpoint/router owns its own GraphQL operation; a shared GraphQLEndpoint
or central Core GraphQL router is not permitted.

Run `make gql-generate` after changing SDL. Generated output must not be edited
by hand.
