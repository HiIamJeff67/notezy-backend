# DurableJob event contracts v1

This package owns the Kafka protocols coordinated by DurableJob:

- RoutineTask claim, assignment, completion, and failure messages.
- Yjs maintenance requests and results exchanged with Core.

Core publishes Yjs maintenance hints from `contracts/core/v1/events`. YjsWorker
owns the maintenance operation, command, and worker-result contracts in
`contracts/yjsworker/v1/events`.

The generic envelope is imported from `contracts/types/event.go`; this package owns the
topics, event types, and payloads. Consumer groups remain runtime deployment
configuration and are defined by the DurableJob/Core transport composition
roots. This package does not import Core repositories, database schemas, or
Kafka clients.
