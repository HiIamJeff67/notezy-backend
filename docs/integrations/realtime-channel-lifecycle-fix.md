# Realtime BlockPack Channel Lifecycle Fix

## Purpose

This document records the frontend change required for the repeated
`channel_not_found` error when opening a BlockPack.

The backend assigns `connectorChannelId` only in the `subscribed` frame. The
value is local to the current physical WebSocket connection and is invalid
after reconnect, resync, or server-side channel detachment.

## Observed failure

The frontend trace has this sequence:

```text
send subscribe
send unsubscribe { connectorChannelId: 2 }
receive error { code: "channel_not_found", connectorChannelId: 2 }
```

The same sequence repeats with `connectorChannelId: 3`. There is no evidence
that the frontend received a new `subscribed` frame before sending the
unsubscribe. The server has already removed the logical channel when it sends
`resubscribe_required`, `permission_revoked`, `resource_unavailable`, or
`channel_backpressure`; sending a second unsubscribe for that ID is invalid.

`dropped awareness update before channel connect` is expected while the
channel is not subscribed. It is not the source of `channel_not_found`.

## Required state rules

Maintain these invariants in `shared/api/websocket/client.ts`:

1. `connectorChannelId` starts as `null`.
2. It is assigned only from a backend `subscribed` frame.
3. It is cleared on reconnect, server-side lifecycle errors, and successful
   unsubscribe.
4. A value from an earlier WebSocket connection must never be reused.
5. A pending subscribe is not unsubscribed by sending a guessed or stale
   connector ID.

## Client changes

### 1. Do not unsubscribe a server-detached channel

In `handleServerError`, the current lifecycle-error branch sends an
`unsubscribe` frame for every error that requires resync. Remove that send for
errors where the backend already detached the channel:

```ts
const serverDetachedChannelCodes = new Set([
  "permission_revoked",
  "resource_unavailable",
  "resync_required",
  "channel_backpressure",
  "worker_unavailable",
  "block_pack_quota_exceeded",
]);
```

For these codes, clear `channelByConnectorId`, set the channel's
`connectorChannelId` and `pendingRequestId` to `null`, stop the Yjs provider,
and wait for an explicit resync operation. Do not send another unsubscribe.

`unsubscribe` is appropriate only when the client intentionally releases a
currently confirmed subscription, such as `releaseBlockPackChannel` after the
retention delay.

### 2. Guard `unregisterBlockPackChannel`

Only send an unsubscribe when all of the following are true:

```ts
channel.connectorChannelId !== null &&
channel.pendingRequestId === null &&
this.isSocketOpen()
```

If `pendingRequestId !== null`, cancel the local pending request mapping and
remove the channel locally. Do not send an unsubscribe because the server has
not assigned a channel ID yet.

### 3. Recover after `resync_required`

The recovery sequence must be:

```text
stop/destroy old Yjs provider
clear connectorChannelId
request a fresh one-time BlockPack channel ticket
send subscribe
wait for subscribed
start a new Yjs provider
```

Do not reuse the consumed channel ticket. Do not start the provider or send
Yjs/awareness binary frames before `subscribed` is received.

`block_pack_quota_exceeded` is a distinct recovery path. Preserve the rejected
local content as an editor draft, destroy the polluted local Y.Doc, request a
fresh ticket, and rebuild from the authoritative initial state. Do not merge the
authoritative state into the rejected local Y.Doc because the rejected CRDT
operations would remain present and be sent again. The ticket response and
`subscribed` frame expose `documentQuotaPolicyVersion` and `maximumBlockCount`
for proactive UI feedback, but the Yjs Worker remains authoritative.

The explicit `resyncBlockPackChannel` operation may create a replacement
channel object, but it must preserve the existing retain count and must not
call `unregisterBlockPackChannel` in a way that sends an unsubscribe for a
server-detached ID.

### 4. Keep React cleanup stable

The cleanup of the BlockPack editor effect should release a channel only when
the editor is actually unmounted or the BlockPack/permission changes. A
normal state rerender must not call `unregisterBlockPackChannel` while a ticket
request or subscribe request is pending.

## Expected successful trace

```text
ready
request fresh BlockPack channel ticket
send subscribe
receive subscribed { connectorChannelId: N }
start Yjs provider
send/receive binary Yjs and awareness frames using N
```

If the server sends `resubscribe_required` later:

```text
receive resubscribe_required { connectorChannelId: N }
stop provider
clear N locally
request fresh ticket
send subscribe
receive subscribed { connectorChannelId: M }
start provider using M
```

`M` may differ from `N`; the client must not send any frame using `N` after
the lifecycle error.

## Verification checklist

- No `unsubscribe` is sent for a channel that has only a pending subscribe.
- No `unsubscribe` is sent in response to `resync_required` or another
  server-detached lifecycle error.
- Every Yjs provider starts after a matching `subscribed` frame.
- A reconnect clears all connector IDs before resubscribing.
- A fresh channel ticket is requested for every resubscribe.
- The browser WebSocket trace contains `subscribed` before any binary frame.
- A quota rejection does not enter an automatic resubscribe loop or resend the
  rejected update.
- `channel_not_found` is not converted into a persistent editor error when it
  was caused by cleanup of an already-detached channel.
