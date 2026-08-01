# YjsWorker v1 contracts

YjsWorker is an independently deployable worker boundary. Its HTTP payloads are
versioned here when the TypeScript worker and Go DurableJob runtime need a
stable request/response schema. Binary Yjs protocol types that are shared with
the WebSocket runtime remain in `internal/shared/types`.

The worker owns Yjs document transformation and does not access the Core
database. Persistence and projection are applied by the Core-owned data path.
