# RealtimeGateway v1 contracts

This package is the versioned public HTTP/WebSocket boundary owned by RealtimeGateway.
API Gateway uses it to request ephemeral WebSocket presence snapshots; it never
reads RealtimeGateway Redis directly. RealtimeGateway returns only public user
IDs, channel permissions, and connection counts. Core remains responsible for
requester authorization and profile enrichment.
