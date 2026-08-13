package main

import "path/filepath"

func writeRealtimeArtifacts(root string) {
	base := filepath.Join(root, "contracts", "realtime-gateway", "v1", "public-api")
	exception := map[string]any{"type": "object", "properties": map[string]any{
		"reason": map[string]any{"type": "string"}, "domain": map[string]any{"type": "string"}, "operation": map[string]any{"type": "string"},
		"message": map[string]any{"type": "string"}, "retryable": map[string]any{"type": "boolean"},
	}}
	participant := map[string]any{"type": "object", "required": []string{"userPublicId", "channelPermission", "connectionCount"}, "properties": map[string]any{
		"userPublicId": map[string]any{"type": "string", "format": "uuid"}, "channelPermission": map[string]any{"type": "string", "enum": []string{"read", "write"}},
		"connectionCount": map[string]any{"type": "integer", "minimum": 1},
	}}
	response := map[string]any{"type": "object", "required": []string{"version", "metadata", "data"}, "properties": map[string]any{
		"version":   map[string]any{"type": "string", "const": "v1"},
		"metadata":  map[string]any{"type": "object", "required": []string{"requestId", "respondedAt"}, "properties": map[string]any{"requestId": map[string]any{"type": "string"}, "respondedAt": map[string]any{"type": "string", "format": "date-time"}}},
		"data":      map[string]any{"type": "object", "required": []string{"participants"}, "properties": map[string]any{"participants": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Participant"}}}},
		"exception": map[string]any{"$ref": "#/components/schemas/Exception"},
	}}
	openAPI := map[string]any{
		"openapi": "3.1.0", "info": map[string]any{"title": "Notezy RealtimeGateway HTTP API", "version": "1.0.0", "description": "Versioned HTTP surface emitted directly by RealtimeGateway. The WebSocket protocol is defined by asyncapi/asyncapi.json."},
		"servers": []any{map[string]any{"url": "http://localhost/realtime/development/v1"}, map[string]any{"url": "https://api.notezy.app/realtime/development/v1"}},
		"paths": map[string]any{
			"/block-pack/{block-pack-id}/participants": map[string]any{"get": map[string]any{
				"operationId": "getRealtimeBlockPackParticipants", "summary": "Get active BlockPack participants", "tags": []string{"Presence"},
				"security":   []any{map[string]any{"accessCookie": []any{}}, map[string]any{"refreshCookie": []any{}}},
				"parameters": []any{map[string]any{"name": "block-pack-id", "in": "path", "required": true, "schema": map[string]any{"type": "string", "format": "uuid"}}, map[string]any{"name": "X-Request-Id", "in": "header", "required": false, "schema": map[string]any{"type": "string"}}},
				"responses":  map[string]any{"200": map[string]any{"description": "Presence snapshot", "content": map[string]any{"application/json": map[string]any{"schema": map[string]any{"$ref": "#/components/schemas/ParticipantsResponse"}}}}, "400": map[string]any{"description": "Invalid request"}, "401": map[string]any{"description": "Authentication failed"}, "429": map[string]any{"description": "Rate limited"}, "503": map[string]any{"description": "Presence unavailable"}},
			}},
		},
		"components": map[string]any{"schemas": map[string]any{"Exception": exception, "Participant": participant, "ParticipantsResponse": response}, "securitySchemes": map[string]any{
			"accessCookie": map[string]any{"type": "apiKey", "in": "cookie", "name": "accessToken"}, "refreshCookie": map[string]any{"type": "apiKey", "in": "cookie", "name": "refreshToken"},
		}},
	}
	writeJSON(filepath.Join(base, "openapi", "openapi.json"), openAPI)
	writeJSON(filepath.Join(base, "asyncapi", "asyncapi.json"), realtimeAsyncAPI(participant))
	writeJSON(filepath.Join(base, "postman", "Notezy-RealtimeGateway-v1.postman_collection.json"), realtimePostman())
	writeJSON(filepath.Join(base, "postman", "Notezy-RealtimeGateway-v1.postman_environment.example.json"), postmanEnvironment())
	writeText(filepath.Join(base, "examples", "control-frames.json"), realtimeControlExamples())
	writeText(filepath.Join(base, "examples", "websocket-browser.js"), realtimeBrowserExample())
	writeRealtimeRules(base)
}

func realtimeAsyncAPI(participant map[string]any) map[string]any {
	return map[string]any{
		"asyncapi": "3.0.0", "info": map[string]any{"title": "Notezy RealtimeGateway WebSocket API", "version": "1.0.0", "description": "Public multiplexed RealtimeGateway protocol. Text frames are JSON; binary frames use the documented six-byte header."},
		"servers": map[string]any{
			"local":      map[string]any{"host": "localhost", "pathname": "/realtime/development/v1", "protocol": "ws", "description": "Local development", "security": []any{map[string]any{"$ref": "#/components/securitySchemes/connectionTicket"}}},
			"hostedBeta": map[string]any{"host": "api.notezy.app", "pathname": "/realtime/development/v1", "protocol": "wss", "description": "Hosted Beta", "security": []any{map[string]any{"$ref": "#/components/securitySchemes/connectionTicket"}}},
		},
		"channels": map[string]any{"realtime": map[string]any{
			"address": "/realtime/development/v1",
			"messages": map[string]any{
				"clientControl": map[string]any{"$ref": "#/components/messages/ClientControl"},
				"serverControl": map[string]any{"$ref": "#/components/messages/ServerControl"},
				"binary":        map[string]any{"$ref": "#/components/messages/BinaryFrame"},
			},
		}},
		"operations": map[string]any{
			"sendClientControl":    map[string]any{"action": "send", "channel": map[string]any{"$ref": "#/channels/realtime"}, "messages": []any{map[string]any{"$ref": "#/channels/realtime/messages/clientControl"}, map[string]any{"$ref": "#/channels/realtime/messages/binary"}}},
			"receiveServerControl": map[string]any{"action": "receive", "channel": map[string]any{"$ref": "#/channels/realtime"}, "messages": []any{map[string]any{"$ref": "#/channels/realtime/messages/serverControl"}, map[string]any{"$ref": "#/channels/realtime/messages/binary"}}},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{"connectionTicket": map[string]any{"type": "httpApiKey", "in": "header", "name": "Sec-WebSocket-Protocol", "description": "Single-use connection ticket issued by Gateway v1."}},
			"messages": map[string]any{
				"ClientControl": map[string]any{"name": "ClientControl", "contentType": "application/json", "payload": map[string]any{"oneOf": []any{
					frameSchema("subscribe", map[string]any{"channelType": map[string]any{"type": "string", "const": "BlockPack"}, "channelId": map[string]any{"type": "string", "format": "uuid"}, "channelTicket": map[string]any{"type": "string"}}, []string{"channelType", "channelId", "channelTicket"}),
					frameSchema("unsubscribe", map[string]any{"connectorChannelId": map[string]any{"type": "integer", "minimum": 1}}, []string{"connectorChannelId"}),
					frameSchema("ack", map[string]any{"connectorChannelId": map[string]any{"type": "integer", "minimum": 1}, "sequence": map[string]any{"type": "integer", "format": "int64"}}, []string{"connectorChannelId", "sequence"}),
					frameSchema("ping", nil, nil), frameSchema("heartbeat", nil, nil),
				}}},
				"ServerControl": map[string]any{"name": "ServerControl", "contentType": "application/json", "payload": map[string]any{"oneOf": serverControlSchemas(participant)}},
				"BinaryFrame":   map[string]any{"name": "BinaryFrame", "contentType": "application/octet-stream", "payload": map[string]any{"type": "string", "format": "binary", "description": "Bytes 0..5 are version, type, and big-endian connectorChannelId; remaining bytes are raw Yjs or awareness payload."}},
			},
		},
	}
}

func serverControlSchemas(participant map[string]any) []any {
	uuid := func() map[string]any { return map[string]any{"type": "string", "format": "uuid"} }
	channelFields := map[string]any{
		"channelType":        map[string]any{"type": "string", "const": "BlockPack"},
		"channelId":          uuid(),
		"connectorChannelId": map[string]any{"type": "integer", "minimum": 1},
	}
	copyFields := func(source map[string]any) map[string]any {
		result := map[string]any{}
		for name, value := range source {
			result[name] = value
		}
		return result
	}
	subscribed := copyFields(channelFields)
	subscribed["existing"] = map[string]any{"type": "boolean"}
	subscribed["documentQuotaPolicyVersion"] = map[string]any{"type": "integer", "const": 1}
	subscribed["maximumBlockCount"] = map[string]any{"type": "integer", "minimum": 1}
	subscribed["participants"] = map[string]any{"type": "array", "items": participant}
	presence := map[string]any{"channelType": channelFields["channelType"], "channelId": channelFields["channelId"], "participant": participant}
	errorCodes := []string{
		"authentication_managed_by_upgrade", "binary_channel_not_ready", "block_pack_quota_exceeded", "channel_backpressure", "channel_limit_exceeded", "channel_not_found", "channel_permission_denied",
		"invalid_acknowledgement", "invalid_binary_frame", "invalid_channel_id", "invalid_channel_ticket", "invalid_channel_type", "invalid_connector_channel_id", "invalid_control_frame",
		"permission_revoked", "resource_unavailable", "room_admission_unavailable", "room_connection_limit_exceeded", "resubscribe_required", "ticket_already_used", "unsupported_binary_type",
		"unsupported_channel_type", "unsupported_control_type", "unsupported_message_type", "unsupported_protocol_version", "worker_unavailable",
	}
	errorFields := copyFields(channelFields)
	errorFields["code"] = map[string]any{"type": "string", "enum": errorCodes}
	errorFields["message"] = map[string]any{"type": "string"}
	return []any{
		frameSchema("ready", map[string]any{"connectionId": uuid(), "resubscribeRequired": map[string]any{"type": "boolean"}}, []string{"connectionId", "resubscribeRequired"}),
		frameSchema("subscribed", subscribed, []string{"channelType", "channelId", "connectorChannelId", "existing", "documentQuotaPolicyVersion", "maximumBlockCount"}),
		frameSchema("unsubscribed", channelFields, []string{"channelType", "channelId", "connectorChannelId"}),
		frameSchema("acknowledged", map[string]any{"connectorChannelId": channelFields["connectorChannelId"], "sequence": map[string]any{"type": "integer", "format": "int64"}}, []string{"connectorChannelId", "sequence"}),
		frameSchema("pong", nil, nil),
		frameSchema("heartbeat", map[string]any{"unixMilliNow": map[string]any{"type": "integer", "format": "int64"}}, []string{"unixMilliNow"}),
		frameSchema("presence-joined", presence, []string{"channelType", "channelId", "participant"}),
		frameSchema("presence-left", presence, []string{"channelType", "channelId", "participant"}),
		frameSchema("presence-updated", presence, []string{"channelType", "channelId", "participant"}),
		frameSchema("resource-event", map[string]any{"eventId": uuid(), "eventType": map[string]any{"type": "string"}, "resourceId": uuid(), "targetUserPublicId": uuid(), "change": map[string]any{"type": "string"}, "permission": map[string]any{"type": "string"}}, []string{"eventId", "eventType", "resourceId", "change"}),
		frameSchema("notification", map[string]any{"eventId": uuid(), "notificationId": uuid(), "notificationType": map[string]any{"type": "string"}, "priority": map[string]any{"type": "string"}, "templateKey": map[string]any{"type": "string"}, "templateVersion": map[string]any{"type": "integer"}, "payload": map[string]any{}, "createdAt": map[string]any{"type": "string", "format": "date-time"}, "expiresAt": map[string]any{"type": []string{"string", "null"}, "format": "date-time"}}, []string{"eventId", "notificationId", "notificationType", "priority", "templateKey", "templateVersion", "payload", "createdAt"}),
		frameSchema("routine-task-lifecycle", map[string]any{"eventId": uuid(), "routineTaskId": uuid(), "routineTaskRecordId": uuid(), "routineId": uuid(), "purpose": map[string]any{"type": "string"}, "status": map[string]any{"type": "string", "enum": []string{"running", "completed"}}, "attempt": map[string]any{"type": "integer", "minimum": 1}, "occurredAt": map[string]any{"type": "string", "format": "date-time"}}, []string{"eventId", "routineTaskId", "routineTaskRecordId", "routineId", "purpose", "status", "attempt", "occurredAt"}),
		frameSchema("error", errorFields, []string{"code", "message"}),
	}
}

func frameSchema(frameType string, extra map[string]any, required []string) map[string]any {
	properties := map[string]any{"version": map[string]any{"type": "integer", "const": 1}, "type": map[string]any{"type": "string", "const": frameType}, "requestId": map[string]any{"type": "string"}}
	for name, schema := range extra {
		properties[name] = schema
	}
	allRequired := append([]string{"version", "type"}, required...)
	return map[string]any{"type": "object", "required": allRequired, "properties": properties}
}

func realtimePostman() map[string]any {
	return map[string]any{
		"info": map[string]any{"_postman_id": "b6353bbc-7456-49f5-b94d-a6948906cb98", "name": "Notezy RealtimeGateway v1", "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json", "description": "RealtimeGateway HTTP presence request. Use the included browser WebSocket example for the multiplex protocol."},
		"item": []any{map[string]any{
			"name":    "Get active BlockPack participants",
			"request": map[string]any{"method": "GET", "header": []any{map[string]any{"key": "User-Agent", "value": "{{userAgent}}"}, map[string]any{"key": "X-Request-Id", "value": "postman-presence-1"}}, "url": map[string]any{"raw": "{{realtimeBaseUrl}}/block-pack/{{blockPackId}}/participants", "host": []string{"{{realtimeBaseUrl}}"}}},
			"event":   []any{map[string]any{"listen": "test", "script": map[string]any{"type": "text/javascript", "exec": []string{"pm.test('Presence response is successful', function () { pm.response.to.have.status(200); });", "pm.test('Participants is an array', function () { pm.expect(pm.response.json().data.participants).to.be.an('array'); });"}}}},
		}},
	}
}

func realtimeControlExamples() string {
	return `{
  "clientToServer": [
    {"version":1,"type":"subscribe","requestId":"sub-1","channelType":"BlockPack","channelId":"00000000-0000-4000-8000-000000000001","channelTicket":"<single-use-channel-ticket>"},
    {"version":1,"type":"ack","requestId":"ack-1","connectorChannelId":1,"sequence":42},
    {"version":1,"type":"ping","requestId":"ping-1"},
    {"version":1,"type":"heartbeat","requestId":"heartbeat-1"},
    {"version":1,"type":"unsubscribe","requestId":"unsub-1","connectorChannelId":1}
  ],
  "serverToClient": [
    {"version":1,"type":"ready","connectionId":"00000000-0000-4000-8000-000000000002","resubscribeRequired":true},
    {"version":1,"type":"subscribed","requestId":"sub-1","channelType":"BlockPack","channelId":"00000000-0000-4000-8000-000000000001","connectorChannelId":1,"existing":false,"documentQuotaPolicyVersion":1,"maximumBlockCount":1000,"participants":[]},
    {"version":1,"type":"acknowledged","requestId":"ack-1","connectorChannelId":1,"sequence":42},
    {"version":1,"type":"pong","requestId":"ping-1"},
    {"version":1,"type":"heartbeat","requestId":"heartbeat-1","unixMilliNow":1767225600000},
    {"version":1,"type":"presence-joined","channelType":"BlockPack","channelId":"00000000-0000-4000-8000-000000000001","participant":{"userPublicId":"00000000-0000-4000-8000-000000000003","channelPermission":"write","connectionCount":1}},
    {"version":1,"type":"presence-updated","channelType":"BlockPack","channelId":"00000000-0000-4000-8000-000000000001","participant":{"userPublicId":"00000000-0000-4000-8000-000000000003","channelPermission":"read","connectionCount":2}},
    {"version":1,"type":"presence-left","channelType":"BlockPack","channelId":"00000000-0000-4000-8000-000000000001","participant":{"userPublicId":"00000000-0000-4000-8000-000000000003","channelPermission":"read","connectionCount":0}},
    {"version":1,"type":"resource-event","eventId":"00000000-0000-4000-8000-000000000004","eventType":"RootShelfPermissionChanged","resourceId":"00000000-0000-4000-8000-000000000005","targetUserPublicId":"00000000-0000-4000-8000-000000000003","change":"permission_updated","permission":"write"},
    {"version":1,"type":"notification","eventId":"00000000-0000-4000-8000-000000000006","notificationId":"00000000-0000-4000-8000-000000000007","notificationType":"security-alert","priority":"important","templateKey":"security.alert","templateVersion":1,"payload":{},"createdAt":"2026-01-01T00:00:00Z"},
	{"version":1,"type":"routine-task-lifecycle","eventId":"00000000-0000-4000-8000-000000000008","routineTaskId":"00000000-0000-4000-8000-000000000009","routineTaskRecordId":"00000000-0000-4000-8000-000000000010","routineId":"00000000-0000-4000-8000-000000000011","purpose":"CreateBlockPack","status":"running","attempt":1,"occurredAt":"2026-01-01T00:00:00Z"},
    {"version":1,"type":"unsubscribed","requestId":"unsub-1","channelType":"BlockPack","channelId":"00000000-0000-4000-8000-000000000001","connectorChannelId":1},
    {"version":1,"type":"error","requestId":"sub-1","connectorChannelId":1,"code":"invalid_channel_ticket","message":"channel ticket is invalid"}
  ],
  "binaryFrameLayoutExamples": {
    "yjsDocumentHex": "010100000001<payload>",
    "awarenessHex": "010200000001<payload>",
    "note": "The first byte is protocol version, second byte is frame type, next four bytes are big-endian connectorChannelId, and the remainder is raw binary payload."
  }
}`
}

func realtimeBrowserExample() string {
	return `// Paste in a browser console after obtaining tickets from Gateway v1.
// Set these values locally; never commit a real ticket.
const realtimeEndpoint = "ws://localhost/realtime/development/v1";
const connectionTicket = "<single-use-connection-ticket>";
const channelTicket = "<single-use-channel-ticket>";
const blockPackId = "00000000-0000-4000-8000-000000000001";

const socket = new WebSocket(realtimeEndpoint, connectionTicket);
socket.binaryType = "arraybuffer";
socket.addEventListener("message", ({ data }) => {
  if (typeof data === "string") {
    const frame = JSON.parse(data);
    console.log("control", frame);
    if (frame.type === "ready") {
      socket.send(JSON.stringify({ version: 1, type: "subscribe", requestId: "sub-1", channelType: "BlockPack", channelId: blockPackId, channelTicket }));
    }
    return;
  }
  const bytes = new Uint8Array(data);
  const connectorChannelId = new DataView(data).getUint32(2, false);
  console.log("binary", { version: bytes[0], type: bytes[1], connectorChannelId, payload: bytes.slice(6) });
});
socket.addEventListener("close", ({ code, reason }) => console.log("closed", code, reason));
socket.addEventListener("error", (error) => console.error("websocket error", error));`
}

func writeRealtimeRules(base string) {
	writeText(filepath.Join(base, "README.md"), `# Notezy RealtimeGateway v1 public API

This directory specifies every versioned endpoint emitted directly by RealtimeGateway v1.

- `+"`openapi/openapi.json`"+`: participant HTTP API (OpenAPI 3.1).
- `+"`asyncapi/asyncapi.json`"+`: public WebSocket messages (AsyncAPI 3.0).
- `+"`rules/`"+`: admission, frame, presence, limit, and failure rules.
- `+"`examples/`"+`: control-frame catalog and a browser WebSocket client.
- `+"`postman/`"+`: importable participant request and credential-free environment.
- `+"`versions/`"+`: development log and cross-version comparison.

Connection and BlockPack channel tickets are issued by the Gateway v1 endpoints in the Gateway OpenAPI contract.`)
	writeText(filepath.Join(base, "rules", "connection-and-admission.md"), `# Connection and admission rules

- WebSocket path: `+"`/realtime/development/v1`"+`.
- Request a connection ticket from `+"`POST /api/development/v1/realtime/connection/ticket`"+`.
- Pass the connection ticket as the sole `+"`Sec-WebSocket-Protocol`"+` value.
- Request a separate BlockPack ticket from `+"`POST /api/development/v1/realtime/channel/block-pack/ticket`"+` for every logical channel.
- Tickets are EdDSA-signed JWTs, expire after five minutes, bind to the User-Agent hash, and are single-use.
- A failed connection or subscription needs a newly issued ticket.
- One physical socket may hold at most 64 active channels. Released connector channel IDs are not reused on that socket.
- A new `+"`ready`"+` frame is a reconnect boundary. Subscribe every required channel again.
- Only BlockPack channels exist in protocol v1.
- Read tickets allow awareness traffic; only write tickets allow Yjs document mutations.

Room admission uses policy version 1 and `+"`reject-new-subscriber`"+`. Subscriber capacity is distributed and includes both read and write subscriptions.`)
	writeText(filepath.Join(base, "rules", "frames.md"), `# Frame rules

Text control frames are UTF-8 JSON and always contain integer `+"`version: 1`"+` and a string `+"`type`"+`. Client requests should carry a unique `+"`requestId`"+` for correlation.

Client control types are `+"`subscribe`"+`, `+"`unsubscribe`"+`, `+"`ack`"+`, `+"`ping`"+`, and `+"`heartbeat`"+`. Authentication is managed by the HTTP upgrade; an `+"`authenticate`"+` frame is rejected.

Binary frame layout:

| Offset | Length | Meaning |
| --- | --- | --- |
| 0 | 1 | Protocol version (`+"`1`"+`) |
| 1 | 1 | `+"`1`"+` Yjs document or `+"`2`"+` awareness |
| 2 | 4 | Unsigned big-endian connectorChannelId |
| 6 | remaining | Raw Yjs/awareness payload |

Binary payloads are never Base64 JSON. A client must wait for `+"`subscribed`"+` before sending them. ACK sequence values never move backwards.`)
	writeText(filepath.Join(base, "rules", "presence-and-events.md"), `# Presence and event rules

`+"`GET /realtime/development/v1/block-pack/{blockPackId}/participants`"+` returns an ephemeral Redis lease snapshot. It is not an authorization source.

Each participant contains only public user ID, read/write channel permission, and active connection count. An empty list means no active lease was observed. Profile data is intentionally excluded.

After subscription, the socket may emit `+"`presence-joined`"+`, `+"`presence-left`"+`, and `+"`presence-updated`"+`. Apply these deltas idempotently. A left participant has connectionCount zero.

`+"`resource-event`"+` is an invalidation hint, not a resource snapshot. Deduplicate with eventId and refetch canonical REST/GraphQL state. Historical events are not replayed after reconnect. User notifications may also arrive on the root connection and must be treated as transient delivery.

`+"`routine-task-lifecycle`"+` is a user-targeted transient execution hint. It reports `+"`running`"+` when DurableJob begins a RoutineTask handler and `+"`completed`"+` only after Core commits the corresponding result. Deduplicate with eventId; after reconnect or when durable state matters, refetch the canonical RoutineTask and RoutineTaskRecord through the normal API.`)
	writeText(filepath.Join(base, "rules", "limits-and-errors.md"), `# Limits and errors

- Maximum public message: 1 MiB.
- Maximum active channels per connection: 64.
- Maximum concurrent root connections per user: 8.
- WebSocket upgrades are limited to 60 per IP per minute, with a 5/second token bucket and burst 10.
- Participant HTTP requests pass the same IP limit and an additional 60-per-user-per-minute limit.
- Per-channel outbound queue: 256 binary frames and 4 MiB payload.
- Awareness is replaceable ephemeral state; queued Yjs updates are never silently dropped.
- Full Yjs queues detach only the affected channel with `+"`channel_backpressure`"+`.
- A failed read or a write blocked for 10 seconds closes the physical socket.
- Server ping/pong and lease heartbeats maintain liveness; reconnect with backoff.

Stable error codes:

`+"`authentication_managed_by_upgrade`"+`, `+"`binary_channel_not_ready`"+`, `+"`block_pack_quota_exceeded`"+`, `+"`channel_backpressure`"+`, `+"`channel_limit_exceeded`"+`, `+"`channel_not_found`"+`, `+"`channel_permission_denied`"+`, `+"`invalid_acknowledgement`"+`, `+"`invalid_binary_frame`"+`, `+"`invalid_channel_id`"+`, `+"`invalid_channel_ticket`"+`, `+"`invalid_channel_type`"+`, `+"`invalid_connector_channel_id`"+`, `+"`invalid_control_frame`"+`, `+"`permission_revoked`"+`, `+"`resource_unavailable`"+`, `+"`room_admission_unavailable`"+`, `+"`room_connection_limit_exceeded`"+`, `+"`resubscribe_required`"+`, `+"`ticket_already_used`"+`, `+"`unsupported_binary_type`"+`, `+"`unsupported_channel_type`"+`, `+"`unsupported_control_type`"+`, `+"`unsupported_message_type`"+`, `+"`unsupported_protocol_version`"+`, and `+"`worker_unavailable`"+`.

On quota rejection, preserve the local draft separately and rebuild from authoritative state. On permission revocation or resource unavailability, stop the provider and do not reuse that connector channel.`)
	writeText(filepath.Join(base, "versions", "dev-log.md"), `# RealtimeGateway v1 development log

## Current contract baseline

- HTTP surface: BlockPack participant presence lookup.
- WebSocket surface: multiplexed BlockPack subscriptions, JSON control frames, binary Yjs and awareness frames, presence, resource events, notifications, and RoutineTask lifecycle hints.
- Contract formats: OpenAPI 3.1 for HTTP and AsyncAPI 3.0 for WebSocket.
- Admission: single-use Gateway-issued connection and channel tickets.
- Stability: Beta namespace with protocol version 1.

Future entries must call out control-frame, binary-layout, limit, admission, close-code, and reconnect changes.`)
	writeText(filepath.Join(base, "versions", "comparison.md"), `# RealtimeGateway API version comparison

Only RealtimeGateway protocol v1 is currently published, so there is no second public protocol to compare yet.

| Capability | v1 Beta |
| --- | --- |
| HTTP contract | OpenAPI 3.1 |
| WebSocket contract | AsyncAPI 3.0 |
| Logical channels | BlockPack |
| Binary protocol header | 6 bytes |
| Connection authentication | Single-use connection ticket |
| Channel admission | Single-use BlockPack ticket |

When another protocol version is introduced, this file must compare upgrade paths, control and binary frames, ticket claims, limits, error codes, reconnect semantics, and compatibility windows.`)
}
