# DurableJob event contracts v1

This package owns the Kafka protocols coordinated by DurableJob:

- RoutineTask claim, assignment, completion, and failure messages.
- Yjs maintenance hints, requests, commands, and results.

The generic envelope is imported from `contracts/events`; this package owns the
topics, consumer groups, event types, and payloads. It does not import Core
repositories, database schemas, or Kafka clients.
