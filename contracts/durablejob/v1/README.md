# DurableJob v1 contracts

This directory is the versioned boundary owned by the DurableJob service. A
caller uses these contracts when it invokes DurableJob; the path does not encode
the caller or transport direction.

Routine task scheduling uses the following Kafka directions:

```text
Core <-> DurableJob: CoreDurableJobRoutineTaskTopic

The event type distinguishes claim requests, assignments, completed results,
and failed results on this single topic.
```

`ClaimRoutineTasksRequestDto` is a capacity request, not a request for one
specific task. Core owns task claiming, task records, scheduling state, and the
transactional outbox. It responds with one `RoutineTaskAssignment` per task
claimed within the requested batch size.

The service owns its runtime, handlers, validation, and execution state. It does
not own Core database schemas or repositories. DurableJob publishes execution
results after its handlers finish. Core consumes those results and is the only
runtime that finalizes RoutineTask and RoutineTaskRecord state. The InboxEvent
table makes result consumption idempotent.
