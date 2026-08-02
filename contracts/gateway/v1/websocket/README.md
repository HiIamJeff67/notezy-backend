# Gateway WebSocket-to-Core v1 contracts

This package is the versioned private boundary from the Gateway's standalone
WebSocket runtime to Core. It owns DTOs for channel permission checks, durable
Yjs document reads/writes, compaction, and block projection.

The WebSocket runtime authenticates each private HTTP call with a component
delegation token (`actor: websocket`). Core validates that token before it maps
the envelope to its service method. Browser connection and channel tickets stay
on the public WebSocket protocol and are not represented here.

The Go-to-YjsWorker multiplexed binary frames remain in `internal/shared/types`:
they are a transport protocol, not a Core business DTO. YjsWorker does not call
this contract and never accesses Core's database.
