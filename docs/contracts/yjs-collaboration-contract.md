# Yjs Collaboration Contract

本文件定義 Notezy 的 BlockPack collaboration document contract。它是前端 BlockNote、public Realtime gateway、Node Yjs worker、Go persistence 與 Block projection 的共同邊界；任何一方不得自行改變 room、fragment、schema 或 sequence 語意。

## Document Identity

| 項目 | 固定值 |
| --- | --- |
| channel type | `BlockPack` |
| channel id | BlockPack UUID |
| room name | `block-pack:{blockPackId}` |
| Y.XmlFragment name | `document-store` |
| document schema id | `notezy.blocknote` |
| initial document schema version | `1` |

Go constants: `YjsBlockPackRoomPrefix`、`YjsBlockPackFragmentName`、`YjsBlockPackSchemaId`、`YjsBlockPackSchemaVersion`。

`document-store` 必須顯式傳給 BlockNote collaboration configuration，例如 `doc.getXmlFragment("document-store")`。不得依賴 BlockNote Yjs utility 的預設 fragment name。

一個 BlockPack 對應一份 logical Yjs document。`Y.Doc` 是 Node worker 在 active room 的記憶體 runtime object；它不是資料庫 entity，也不會直接傳給前端或 Go service。

## BlockNote Schema

schema version `1` 的 block type manifest 與目前後端 `BlockType` 對齊：

`paragraph`、`heading`、`quote`、`bulletListItem`、`numberedListItem`、`checkListItem`、`toggleListItem`、`image`、`video`、`audio`、`file`、`table`、`codeBlock`。

前端必須以單一 `BlockNoteSchema` factory 建立 editor、Yjs import/export 與 server-side projector 使用的 schema。Node worker 使用相同的 block/inline/style manifest；Go 不解析 Yjs tree，也不自行重建 BlockNote document。

Node projector 使用 `@blocknote/core/yjs` 的 `yXmlFragmentToBlocks`，並明確讀取 `document-store` fragment。它的 schema 排除 `divider`，因為目前後端 `BlockType` 未支援此 block type；前端 schema 也不得建立 `divider`。

新增、刪除或變更 block props、inline content、style schema 都是 schema migration，不是一般 feature flag。

## Version Policy

`schemaVersion` 是 logical document version，不等同 npm package version，也不等同 Realtime protocol version。

| 規則 | 語意 |
| --- | --- |
| 新 document | 以目前 supported schema version 建立，初始為 `1`。 |
| 讀取 | client 與 Node worker 僅可開啟自己明確支援的 version。 |
| 向後相容變更 | 保持 version，僅限舊 reader 能無損理解的新增 optional data。 |
| 不相容變更 | 建立新 version，Node worker 對完整 Y.Doc migration，產出新 snapshot 後才切換。 |
| 投影 | projector 的 schema version 必須與 source document version 相同。 |

目前所有 document 都使用 `YjsBlockPackSchemaVersion = 1`。第一個需要同時支援多個 document schema version 的 migration，才新增 per-document `SchemaVersion`；在那之前不得預先擴充 `Block` schema。

## Persistence And Sequence

durable Yjs truth 是 `BlockPackYjsDocument.Snapshot` 加上尚未 compact 的 `BlockPackYjsUpdate` tail。Snapshot 是 Yjs encoded state update，`StateVector` 是同一個 snapshot 的 encoded state vector；active `Y.Doc` 只是這份 durable truth 的 memory materialization。

每個 BlockPack 必須在建立它的同一筆 transaction 內建立唯一的 `BlockPackYjsDocument`；讀取、append 與 projection 路徑不得 lazy create document。

`BlockTable` 是 Yjs document 的 materialized projection，Block 不支援 soft delete。projection 對不再存在於 document 的 block 使用實體 `DELETE`；BlockPack soft delete 時則保留它的 Blocks，還原 BlockPack 後可直接重用既有 projection。

Block REST read endpoints 與 GraphQL `searchBlocks` 都只讀 `BlockTable` projection；它們不得用於建立或回填 active `Y.Doc`。BlockPack REST read response 會帶 `lastUpdateSequence`、`compactedUntilSequence`、`projectedUntilSequence` 與 `isProjectionCurrent`。前端可以用 `isProjectionCurrent = false` 顯示 read model 正在落後，但 editor 的協作狀態一律仍以 Yjs channel 為準。

| 欄位 | 語意 |
| --- | --- |
| `UpdateSequence` | 單一 BlockPack 內 append-only 的 update 序號，從 `1` 起，永不重用。 |
| `LastUpdateSequence` | 該 BlockPack 已接受的最高 update sequence；不得回退。 |
| `CompactedUntilSequence` | 已被目前 Snapshot 吸收的最高 sequence；不得回退。 |
| `ProjectedUntilSequence` | BlockTable 已成功投影的最高 sequence；document-level checkpoint，初始為 `-1`，且不得回退。 |

不變條件：`0 <= CompactedUntilSequence <= LastUpdateSequence`、`-1 <= ProjectedUntilSequence <= LastUpdateSequence`。空 document 的 durable update/compaction sequence 都是 `0`，未投影時 `ProjectedUntilSequence` 是 `-1`。

compaction 在 Node worker 重建完整 Y.Doc 後執行：它讀取 snapshot 與 update tail、合併到 runtime Y.Doc、寫入新的 Snapshot/StateVector，最後將 `CompactedUntilSequence` 推進到被吸收的最高 sequence。Go 不執行 CRDT merge。

room cold start 固定依序執行：建立空 `Y.Doc`、套用非空 Snapshot、再套用 `CompactedUntilSequence < update_sequence <= LastUpdateSequence` 的 tail。`LastUpdateSequence` 是最新 accepted update，不是 snapshot tail 的查詢起點。

## Public Connection And Capability

root WebSocket authentication 以 connection ticket 取代 access-token middleware；它只識別 user，不授權任何 BlockPack。`NOT-5` 已發出 capability ticket，`NOT-7` 會在每個 BlockPack subscribe 驗證後才建立 channel。

ticket claims 的最小集合：

```json
{
  "sub": "user public UUID",
  "jti": "ticket trace UUID",
  "channelType": "BlockPack",
  "channelId": "blockPack UUID",
  "permission": "read or write",
  "realtimeProtocolVersion": 1,
  "schemaVersion": 1,
  "iat": 0,
  "exp": 0
}
```

`NOT-7` 的 Go Gateway 負責驗證 connection/channel ticket，以及兩者的 user、channel type 與 BlockPack id 是否相符；Node worker 只信任 Go 已驗證並轉送的 attach message。ticket 是短效 stateless capability，`jti` 不代表可在沒有共享 state 下強制一次性使用。permission 被撤銷時，Gateway 送出 `permission_revoked`，移除該 channel，並向 worker 轉送 detach。

## Cross-Service Frames

public WebSocket 的 JSON control frame 與 binary frame header 定義在 `realtime-protocol-contract.md`。Go-to-worker internal binary frame 一律帶有 `connectionId`、`connectorChannelId`、`channelType`、`channelId`；raw Yjs update 不得 Base64 或改寫成 JSON block event。

internal attach/detach 是 idempotent。worker reconnect 後，Gateway 為其所屬 active channels replay attach；worker 會先向 Go cold-load snapshot + tail，materialize `Y.Doc` 後才回傳 complete encoded state。worker 會先套用收到的 raw Yjs updates，再以同一個 BlockPack room 為單位暫存並使用 `Y.mergeUpdates()` 合併為一筆 persistence batch；只有收到 persistence ACK 後才 broadcast merged raw Yjs update。

每個 persistence batch 有只供 Go/worker 使用的 UUID idempotency key。Go 以 `(block_pack_id, persistence_batch_id)` 保證 internal WebSocket retry 不會建立重複 update row；同一 batch 的多個來源 connection 不可任意挑選其中一個寫入 `OriginConnectionId`，必須保留為 `NULL`。append terminal failure 時 worker 對 room 所有 subscriber 發出 `resync_required`，不能 broadcast 未持久化 update。

batch flush 條件由 worker constants 控制：trailing debounce、maximum wait、raw update count、raw payload bytes、最後 subscriber detach 與 graceful worker shutdown。`LastUpdateSequence` 只會在 merged update transaction 成功後推進；每一筆 merged update 只消耗一個 sequence。

## Projection Contract

Node worker 是唯一的 Yjs CRDT merge owner，也是 Y.Doc -> BlockNote blocks conversion owner。它以 current `schemaVersion` 將 active Y.Doc 轉換為 canonical BlockNote block tree，再送出 projection payload 給 Go。

projection payload 最小欄位：

```json
{
  "schemaId": "notezy.blocknote",
  "schemaVersion": 1,
  "projectedSequence": 42,
  "blocks": []
}
```

projection 使用 private Go-to-worker internal frame `apply-block-projection`；BlockPack identity 取自 frame header 的 `channelId`，不是 payload。Go 僅在 payload 的 `schemaId`/`schemaVersion` 受支援、target sequence 不回退，且所屬 BlockPack 有效時寫入 Block projection，成功時回覆 `block-projection-applied`，否則回覆 `block-projection-failed`。`NOT-13` 定義 `blocks` 的 bulk apply、anti-regression transaction 與 accounting semantics；不新增 per-block ordering/search/hash metadata，除非實際 read requirement 證明需要。

worker 僅在 room 沒有尚未持久化的 update 時，將目前 `LastUpdateSequence` 的 document 投影。它對更新 burst 做 debounce，且每個 room 同時最多一筆 projection；只有收到 `block-projection-applied` 後才推進 in-memory `ProjectedUntilSequence`。失敗不會前進 checkpoint，並以 retry delay 重試。

前端不讀取 `BlockPackYjsDocument` 或 `BlockPackYjsUpdate` rows，也不自行合併 update tail；加入 room 時由 Node worker 從 snapshot + tail 恢復 Y.Doc，再以標準 Yjs sync protocol 完成同步。
