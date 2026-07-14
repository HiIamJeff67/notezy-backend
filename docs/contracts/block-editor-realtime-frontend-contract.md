# Block Editor Realtime Frontend Contract

本文件定義前端將既有 `RealtimeWebSocketSmokeProvider` 與 `BlockEditorProvider` 遷移至正式 Notezy Realtime/Yjs 架構時的邊界。它適用於 BlockNote editor；不是 REST block mutation 的替代包裝層。

## Source Of Truth

| 資料 | 用途 | 不可用於 |
| --- | --- | --- |
| active `Y.Doc` | editor 載入、輸入、undo、巢狀、移動、刪除與協作 | REST block write |
| `BlockTable` | REST/GraphQL projection read、搜尋與未來 AI read model | 初始化或回填 active `Y.Doc` |
| `BlockPackYjsDocument` + update log | 僅 Go persistence 與 Node worker 使用 | 前端直接讀取 |

前端不得同時把同一個 editor mutation 寫入 REST 與 Yjs。

## Root Connection

每個 app instance 建立一條 root WebSocket。正式 provider 應改造現有 `RealtimeWebSocketSmokeProvider`，而不是為每個 BlockPack 開一條 socket。

1. provider mount 後，以既有 authenticated REST client 呼叫：

   ```text
   POST /api/development/v1/realtime/createMyRealtimeConnectionTicket
   ```

2. response 提供 `realtimeEndpoint`、`connectionTicket`、`realtimeProtocolVersion`、`expiresAt`。local URL 為：

   ```text
   ws://localhost/realtime/development/v1
   ```

3. 建立 socket 時，`connectionTicket` 是唯一 subprotocol；不可放 query string、Bearer header 或第一個 JSON message：

   ```ts
   const socket = new WebSocket(url, [connectionTicket]);
   ```

4. Gateway 成功驗證後會傳送 `ready`。每個新的 `ready` 都是 reconnect boundary，provider 必須重新取得需要的 channel ticket 並重新 subscribe。

5. provider 的 `useEffect` 必須註冊 `open`、`message`、`close`、`error` listener；cleanup 時取消 reconnect timer、移除 listener 並關閉 socket。connection ticket 只在 upgrade 時使用；每次 reconnect 都要重新取得新的 connection ticket。

`authenticate` control frame 不可使用，Gateway 會回 `authentication_managed_by_upgrade`。

## BlockPack Channel

使用者開啟一個 BlockPack editor 時，`BlockEditorProvider` 必須：

1. 依使用者對 RootShelf 的有效權限請求 `read` 或 `write` channel ticket：

   ```text
   POST /api/development/v1/realtime/createMyBlockPackChannelTicket
   {
     "blockPackId": "UUID",
     "permission": "read" | "write"
   }
   ```

   只讀使用者請求 `read`；可寫使用者請求 `write`。不要把權限不足當成 client-side 可忽略錯誤。

2. 將 response 中的 `channelTicket` 送入 root socket：

   ```json
   {
     "version": 1,
     "type": "subscribe",
     "requestId": "client-generated-id",
     "channelType": "BlockPack",
     "channelId": "BLOCK_PACK_UUID",
     "channelTicket": "ticket"
   }
   ```

3. 收到 `subscribed` 後，保存 `blockPackId -> connectorChannelId` 對應。`connectorChannelId` 是一條 root connection 內的 routing id，不是 BlockPack id，也不能跨 reconnect 重用。

4. editor 關閉時送 `unsubscribe`，移除本地 channel map、Yjs update listener 與 awareness state。root socket 仍保留給其他已開啟的 BlockPack。

## BlockNote And Yjs Provider

channel ticket response 中的 `roomName`、`fragmentName`、`schemaId`、`schemaVersion` 必須先驗證為目前支援值，再建立 BlockNote collaboration configuration。v1 固定使用：

```ts
const fragment = ydoc.getXmlFragment("document-store");
```

自訂 provider 的責任：

- 收到該 `connectorChannelId` 的 `yjs-document` binary payload 時，以明確的 remote origin 呼叫 `Y.applyUpdate`。
- 對本地 `Y.Doc` update，只有 origin 不是 remote origin 時才包裝為該 channel 的 binary frame 送出。
- 對同一份 Yjs update 的 sender echo 可以安全忽略或套用；不得再次送出造成 feedback loop。
- `permission: read` 時不得送 document update；UI 必須是 read-only。awareness 的可用範圍依 realtime protocol contract。
- `permission_revoked` 時立即停止 editor、移除 channel state，並呈現不可再編輯/檢視狀態。
- `resync_required` 時停止該 channel 的即時寫入並重新取得 channel ticket、subscribe；不可用 Block REST rows 補 document。local `Y.Doc` 不得先被銷毀：前端要先以 temporary `Y.Doc` 套用 server complete state、取得 server state vector，再將 local `Y.Doc` 相對於該 vector 缺少的 raw Yjs update 補送，最後才以 server/worker 的正常回應確認同步完成。

public binary header、frame type 與 connector channel routing 以 `realtime-protocol-contract.md` 為準；前端不可自行改寫 raw Yjs update 成 JSON 或 Base64。

## Projection Reads

Block REST read endpoints 僅提供 BlockTable projection：

```text
GET /api/development/v1/block/getMyBlockById
GET /api/development/v1/block/getMyBlocksByIds
GET /api/development/v1/block/getMyBlocksByBlockPackId
```

沒有任何 public Block REST mutation endpoint。BlockPack REST read responses 會提供：

```json
{
  "lastUpdateSequence": 42,
  "compactedUntilSequence": 30,
  "projectedUntilSequence": 42,
  "isProjectionCurrent": true
}
```

`isProjectionCurrent: false` 只代表 projection read model 落後 durable Yjs state；它不是 editor loading source，也不表示應切換回 REST mutation。

## Test Replacement

既有 smoke test 應遷移為正式整合測試，而非 production provider：

- connection ticket 透過 WebSocket subprotocol 成功 upgrade，並收到 `ready`。
- BlockPack channel ticket 可以 subscribe，取得 `connectorChannelId`，並接收完整 Yjs document state。
- local Yjs update 寫入同一 channel，另一個 client 收到並套用。
- read-only channel 無法送 Yjs document update。
- reconnect 後重新取得 ticket、重新 subscribe；`permission_revoked` 與 `resync_required` 會正確 cleanup/recreate provider。

完整 protocol 以 `realtime-protocol-contract.md` 與 `yjs-collaboration-contract.md` 為準。

## Browser Durable Recovery And IndexedDB Storage

- 每個 BlockPack 的 local `Y.Doc` 必須在送出 local update 前持久化至 IndexedDB；可使用 `y-indexeddb` 或等價實作，但不得把 raw Yjs update 寫進 localStorage。
- reconnect、worker restart 或 `resync_required` 後，前端使用 temporary server document 的 state vector 計算 local diff；不得直接丟棄未確認的 local document state。
- 背景圖片 cache 必須和 Yjs document persistence 使用可獨立清理的 IndexedDB object store。背景圖片 metadata 至少保存 byte size 與 last accessed time，總使用量硬上限為 `1 GiB`；寫入前必須先檢查 projected size，並處理 browser `QuotaExceededError`。
- `PreferenceSettingPanel` 必須顯示：local Yjs document cache 的 document count/bytes、背景圖片 cache 的 item count/bytes 與 `1 GiB` 上限、瀏覽器 origin storage usage/quota estimate。提供清除未使用背景圖片、清除全部背景圖片與清除 local document cache 的管理操作；清除 local document cache 前必須提示使用者尚未同步的變更可能遺失。
