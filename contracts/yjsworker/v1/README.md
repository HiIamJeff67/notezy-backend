# YjsWorker v1 contracts

YjsWorker is an independently deployable worker boundary. `contract.go` defines
the versioned Kafka command/reply envelopes; `yjs.go` defines the command data
for loading, appending, compacting, and projecting one BlockPack document.

Commands travel on `notezy.yjsworker.core.command.v1`; Core replies through its
transactional outbox on `notezy.core.yjsworker.reply.v1`. Both use the BlockPack
UUID as their Kafka key, preserving one document's ordering. The worker owns
Yjs transformation but never accesses Core's database directly.

Binary browser and Go-to-YjsWorker socket frames remain an implementation detail
of the realtime transport. They are not persistence contracts.
