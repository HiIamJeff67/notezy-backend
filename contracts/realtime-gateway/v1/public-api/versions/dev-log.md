# RealtimeGateway v1 development log

## Current contract baseline

- HTTP surface: BlockPack participant presence lookup.
- WebSocket surface: multiplexed BlockPack subscriptions, JSON control frames, binary Yjs and awareness frames, presence, resource events, and notifications.
- Contract formats: OpenAPI 3.1 for HTTP and AsyncAPI 3.0 for WebSocket.
- Admission: single-use Gateway-issued connection and channel tickets.
- Stability: Beta namespace with protocol version 1.

Future entries must call out control-frame, binary-layout, limit, admission, close-code, and reconnect changes.
