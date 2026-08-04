# Event contracts v1

This package is Core's versioned domain-event boundary for Kafka. It has no
dependency on Core implementation, RealtimeGateway, persistence, or a Kafka
client.

`notezy.core.lifecycle.v1` carries Core lifecycle facts. Every event uses its
aggregate UUID string as `kafkaKey`, so one aggregate remains ordered within a
Kafka partition. `BlockPack` events always describe exactly one BlockPack;
bulk mutations create one outbox row and one event per affected BlockPack.

`TargetUserPublicId` is optional. A present value targets only that user's
RealtimeGateway channels; an omitted value applies to all channels for the
aggregate. It always contains the public user UUID, never Core's internal user
ID.

The envelope and payloads exclude entity snapshots, browser credentials,
cookies, tokens, presence, awareness, and Yjs updates. Producers may add only
backward-compatible optional fields within `v1`; a breaking change requires a
new contract and topic major version.
