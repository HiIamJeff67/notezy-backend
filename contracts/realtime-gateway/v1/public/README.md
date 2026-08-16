# Notegic RealtimeGateway v1 public contract

This directory specifies every versioned endpoint emitted directly by RealtimeGateway v1 for realtime clients. It is a runtime-specific public contract, separate from the APIGateway integration contract.

- `openapi/openapi.json`: participant HTTP API (OpenAPI 3.1).
- `asyncapi/asyncapi.json`: public WebSocket messages (AsyncAPI 3.0).
- `rules/`: admission, frame, presence, limit, and failure rules.
- `examples/`: control-frame catalog and a browser WebSocket client.
- `postman/`: importable participant request and credential-free environment.
- `versions/`: development log and cross-version comparison.

Connection and BlockPack channel tickets are issued by ClientGateway v1 endpoints. They are not public APIGateway operations.
