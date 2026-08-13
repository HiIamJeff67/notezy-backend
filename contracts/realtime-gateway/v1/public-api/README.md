# Notezy RealtimeGateway v1 public API

This directory specifies every versioned endpoint emitted directly by RealtimeGateway v1.

- `openapi/openapi.json`: participant HTTP API (OpenAPI 3.1).
- `asyncapi/asyncapi.json`: public WebSocket messages (AsyncAPI 3.0).
- `rules/`: admission, frame, presence, limit, and failure rules.
- `examples/`: control-frame catalog and a browser WebSocket client.
- `postman/`: importable participant request and credential-free environment.
- `versions/`: development log and cross-version comparison.

Connection and BlockPack channel tickets are issued by the Gateway v1 endpoints in the Gateway OpenAPI contract.
