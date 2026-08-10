# YjsWorker v1 contracts

YjsWorker is an independently deployable worker boundary. `contract.go` defines
the versioned Kafka command/reply envelopes; `yjs.go` defines the command data
for loading, appending, compacting, and projecting one BlockPack document.

Commands travel on `notezy.adapters.core.command.v1`; Core replies through its
transactional outbox on `notezy.core.adapters.reply.v1`. Both use the BlockPack
UUID as their Kafka key, preserving one document's ordering. The worker owns
Yjs transformation but never accesses Core's database directly.

Binary browser and Go-to-YjsWorker socket frames remain an implementation detail
of the realtime transport. They are not persistence contracts.

Protocol constants are owned by this versioned contract package. The TypeScript
counterpart is colocated at `contracts/yjs-worker/v1/yjsworker_contract.ts` so the
isolated Node runtime and Go consumers use the same versioned definitions. The
planned contract generator will keep the TypeScript file synchronized with the
Go definitions.
