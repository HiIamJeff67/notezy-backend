# RealtimeGateway API version comparison

Only RealtimeGateway protocol v1 is currently published, so there is no second public protocol to compare yet.

| Capability | v1 Beta |
| --- | --- |
| HTTP contract | OpenAPI 3.1 |
| WebSocket contract | AsyncAPI 3.0 |
| Logical channels | BlockPack |
| Binary protocol header | 6 bytes |
| Connection authentication | Single-use connection ticket |
| Channel admission | Single-use BlockPack ticket |

When another protocol version is introduced, this file must compare upgrade paths, control and binary frames, ticket claims, limits, error codes, reconnect semantics, and compatibility windows.
