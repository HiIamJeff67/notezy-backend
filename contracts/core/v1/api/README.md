# Core API contracts

This package contains the Core-owned RequestDto and ResponseDto types used by
its HTTP and GraphQL business operations. Each business domain has a dedicated
folder named after its route domain. Within it, files are named for one
operation family: `get.go`, `create.go`, `update.go`, `restore.go`, `delete.go`,
`permission.go`, `search.go`, or `visualization.go`.

`contract.go` holds only the reusable `RequestDto` composition type. It is not
an HTTP envelope; Gateway owns that envelope in `contracts/gateway/v1`.
