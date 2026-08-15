# Notezy Backend Conventions

This directory is the shared baseline for Codex and developers changing the
backend. The existing codebase is the source of truth. Items marked as
recommendations may be adjusted to match team preferences; update these files
directly after changing a decision so future agents follow the latest rules.

## How to Use

1. Before implementation, read the documents that match the change scope.
2. When an existing pattern conflicts with these conventions, preserve the
   existing pattern first; open a clearly scoped refactoring task if it should
   be unified.
3. Record new cross-module decisions in this directory instead of leaving them
   only in a PR or chat history.

## Document Index

| Document | Scope |
| --- | --- |
| [01-general-go.md](01-general-go.md) | Go formatting, naming, dependency direction, and change scope |
| [02-architecture.md](02-architecture.md) | Gateway/API/worker ownership, HTTP request flow, GraphQL, and background work |
| [03-http-api.md](03-http-api.md) | routes, controllers, adapters, request/response contracts, exceptions, and observability |
| [04-persistence.md](04-persistence.md) | services, repositories, scopes, schemas, transactions, and SQL |
| [05-testing-and-generated-code.md](05-testing-and-generated-code.md) | tests, test data, generated GraphQL code, and verification checklist |
| [06-exceptions.md](06-exceptions.md) | base/service exception domains, error origins, and `exceptions.Cover()` |
| [07-version-control.md](07-version-control.md) | commit messages, daily devlogs, and Git hooks |

## Priority Order

1. Correctness, security, data integrity, and existing APIs/contracts.
2. Explicit conventions in this directory.
3. Established patterns near the same feature.
4. Go conventions and the smallest readable, testable implementation.

Do not add an abstraction layer, interface, or dependency for a possible future
need; add one only when the need exists and the current pattern cannot support
it.

Ownership for the target workspace and staged migration follows
[microservice-architecture.md](../codebase-design/microservice-architecture.md).
Do not create an empty package or temporary wrapper merely to match the target
directory diagram before code has migrated to that path.
