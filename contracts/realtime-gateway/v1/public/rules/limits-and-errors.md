# Limits and errors

- Maximum public message: 1 MiB.
- Maximum active channels per connection: 64.
- Maximum concurrent root connections per user: 8.
- WebSocket upgrades are limited to 60 per IP per minute, with a 5/second token bucket and burst 10.
- Participant HTTP requests pass the same IP limit and an additional 60-per-user-per-minute limit.
- Per-channel outbound queue: 256 binary frames and 4 MiB payload.
- Awareness is replaceable ephemeral state; queued Yjs updates are never silently dropped.
- Full Yjs queues detach only the affected channel with `channel_backpressure`.
- A failed read or a write blocked for 10 seconds closes the physical socket.
- Server ping/pong and lease heartbeats maintain liveness; reconnect with backoff.

Stable error codes:

`authentication_managed_by_upgrade`, `binary_channel_not_ready`, `block_pack_quota_exceeded`, `channel_backpressure`, `channel_limit_exceeded`, `channel_not_found`, `channel_permission_denied`, `invalid_acknowledgement`, `invalid_binary_frame`, `invalid_channel_id`, `invalid_channel_ticket`, `invalid_channel_type`, `invalid_connector_channel_id`, `invalid_control_frame`, `permission_revoked`, `resource_unavailable`, `room_admission_unavailable`, `room_connection_limit_exceeded`, `resubscribe_required`, `ticket_already_used`, `unsupported_binary_type`, `unsupported_channel_type`, `unsupported_control_type`, `unsupported_message_type`, `unsupported_protocol_version`, and `worker_unavailable`.

On quota rejection, preserve the local draft separately and rebuild from authoritative state. On permission revocation or resource unavailability, stop the provider and do not reuse that connector channel.
