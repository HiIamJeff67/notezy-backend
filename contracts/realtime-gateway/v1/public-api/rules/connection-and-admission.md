# Connection and admission rules

- WebSocket path: `/realtime/development/v1`.
- Request a connection ticket from `POST /api/development/v1/realtime/connection/ticket`.
- Pass the connection ticket as the sole `Sec-WebSocket-Protocol` value.
- Request a separate BlockPack ticket from `POST /api/development/v1/realtime/channel/block-pack/ticket` for every logical channel.
- Tickets are EdDSA-signed JWTs, expire after five minutes, bind to the User-Agent hash, and are single-use.
- A failed connection or subscription needs a newly issued ticket.
- One physical socket may hold at most 64 active channels. Released connector channel IDs are not reused on that socket.
- A new `ready` frame is a reconnect boundary. Subscribe every required channel again.
- Only BlockPack channels exist in protocol v1.
- Read tickets allow awareness traffic; only write tickets allow Yjs document mutations.

Room admission uses policy version 1 and `reject-new-subscriber`. Subscriber capacity is distributed and includes both read and write subscriptions.
