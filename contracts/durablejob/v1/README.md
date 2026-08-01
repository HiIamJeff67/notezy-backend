# DurableJob v1 contracts

This directory is the versioned boundary owned by the DurableJob service. A
caller uses these contracts when it invokes DurableJob; the path does not encode
the caller or transport direction.

DurableJob currently receives work through the shared PostgreSQL task tables and
does not expose a public HTTP API. The contract package is reserved for the
internal command and result envelopes that will replace direct task-table
coordination when the Core-to-DurableJob transport is introduced.

The service owns its runtime, handlers, validation, exceptions, and data model
under `internal/services/durablejob`. It may share the existing database during
this migration, but it must not import Core service implementations.
