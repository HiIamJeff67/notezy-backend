# Frame rules

Text control frames are UTF-8 JSON and always contain integer `version: 1` and a string `type`. Client requests should carry a unique `requestId` for correlation.

Client control types are `subscribe`, `unsubscribe`, `ack`, `ping`, and `heartbeat`. Authentication is managed by the HTTP upgrade; an `authenticate` frame is rejected.

Binary frame layout:

| Offset | Length | Meaning |
| --- | --- | --- |
| 0 | 1 | Protocol version (`1`) |
| 1 | 1 | `1` Yjs document or `2` awareness |
| 2 | 4 | Unsigned big-endian connectorChannelId |
| 6 | remaining | Raw Yjs/awareness payload |

Binary payloads are never Base64 JSON. A client must wait for `subscribed` before sending them. ACK sequence values never move backwards.
