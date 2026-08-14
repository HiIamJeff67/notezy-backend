# 前端 Contract 遷移：Gateway 分流與 RoutineTask 月度額度

> 暫時整合文件。本文比較目前工作區與上一個提交，供 Notezy
> frontend 在後端部署這批變更時串接；它不取代正式的 generated contract。

## 結論

現有 Web client 的 REST route、DTO JSON key 與 cookie/JWT 登入流程都沒有
破壞性變更。前端需要調整的是：

1. 將 `routineTaskCostUnitCount` 視為「本月已消耗的執行額度」，不再作為
   建立、編輯或刪除 RoutineTask 時的 payload 額度。
2. 移除建立／更新 RoutineTask 前依估算 payload cost 對 plan limit 的前端
   阻擋，以及相關 optimistic quota 加減。
3. Web client 繼續走 ClientGateway 的 cookie/JWT flow；**不要**把 API key
   放進瀏覽器、localStorage、URL 或前端環境變數。

## Contract 位置變更

| 先前位置 | 現在位置 | 前端用途 |
| --- | --- | --- |
| `contracts/gateway/v1/public-api/` | 已 deprecated，暫留作遷移相容 | 不要再作為文件或 code generation 的來源 |
| — | `contracts/client-gateway/v1/public/` | 現有 Web client 的 JWT cookie、登入、註冊與 client-only flow |
| — | `contracts/api-gateway/v1/public/` | 給外部整合、CLI 或使用者自有 server 的公開 API |
| `contracts/realtime-gateway/v1/public-api/` | `contracts/realtime-gateway/v1/public/` | RealtimeGateway 的 runtime public contract；既有前端 realtime 串接維持既有 ticket flow |

`contracts/api-gateway/v1/public/` 的 OpenAPI、Postman 與 examples
才是唯一公開 API 文件。它的 authenticated request 使用：

```http
X-API-Key: nzy_<secret>
```

這是 server-to-server／CLI integration 的認證方式，不是瀏覽器登入方式。當前
Web client 仍使用 HttpOnly access/refresh cookies、`credentials: "include"` 與
既有 CSRF 行為；不需替現有 API client 統一加上 `X-API-Key`。

## RoutineTask 月度執行額度

`GET /api/development/v1/me/account/` 的 response 沒有新增、刪除或改名欄位：

```ts
routineTaskCostUnitCount: number
```

但它的語意已改變：

| 項目 | 舊行為 | 新行為 |
| --- | --- | --- |
| 值的來源 | 已建立的 RoutineTask payload cost 總和 | 當月實際被 Core claim 並開始執行的 RoutineTask cost 總和 |
| 何時改變 | 建立、更新、刪除 task，或 ownership transfer | Task 被排程 worker claim 執行時 |
| 刪除／編輯 task | 會讓用量回扣或變動 | 不會回扣或變動已消耗的月度用量 |
| 額度歸屬 | 曾受 task 關聯／轉移影響 | 永遠消耗建立該 RoutineTask 的 actor monthly quota |

`RoutineTask.costUnit` 仍表示該 task 的單次執行 cost；它由 payload 大小計算，
但前端不再用它預先拒絕 create/update。後端會在實際 claim 時原子地檢查與扣除。
當額度不足時，task 保持 `Idle`，不會被派送或標記為失敗；後端週期性重設額度後，
它才會再次符合執行資格。

### 月度額度上限

請同步更新任何前端顯示用 fallback constants。後端 `PlanLimitation` 的
`maxRoutineTaskCostUnitCount` 現在是每月執行額度：

| Plan | 每月 RoutineTask cost units |
| --- | ---: |
| Free | 100 |
| Pro | 300 |
| Premium | 600 |
| Ultimate | 1,200 |
| Enterprise | 6,000 |

## 前端必做調整

### 保留 schema，但改變顯示語意

現有 Zod schema 的 `routineTaskCostUnitCount` 可原樣保留。建議在 UI state 中改用
`routineTaskMonthlyCostUnitUsed` 這類 local variable name，文案改為「本月已使用的
RoutineTask 執行額度」，避免再暗示它是已建立模板／payload 的總大小。

### 移除 create/update 的預先阻擋

目前 `CreateRoutineTaskDialog` 與 `RoutineTaskInspector` 有根據 payload estimate、
`routineTaskCostUnitCount` 和 plan limit 阻擋送出的邏輯。請移除該 quota gate；保留
一般欄位、payload 格式與表單 validation 即可。

同時移除以下 optimistic mutation：

- create 後本地增加 account 的 `routineTaskCostUnitCount`；
- update 後以新舊 payload cost 差額調整它；
- delete 或 ownership transfer 後本地回扣／轉移它。

若 UI 想顯示單次預估 cost，可以繼續計算，但文案必須是「此 task 每次執行預估
消耗」；不可用它判斷使用者能否建立或編輯 task。

### 用量更新策略

create/update/delete response 不再代表月度用量已改變。當帳戶頁、upgrade UI 或 quota
indicator 需要最新數字時，重新取得 `GET /me/account/` 即可。這批變更沒有新增一個
用於 quota refresh 的 REST response field 或 WebSocket event，因此不要假設 task
saved 後用量會立刻改變。

## 前端驗收清單

- [ ] 既有 Web API client 保持 cookie/JWT flow，沒有把 API key 寫入 browser code。
- [ ] `routineTaskCostUnitCount` 仍可 decode，UI 文案改為當月已使用的執行額度。
- [ ] RoutineTask create/update 不再因本地 payload-cost estimate 被拒絕。
- [ ] create/update/delete/transfer 不再 optimistic 修改帳戶的 cost unit count。
- [ ] plan limitation 的顯示值更新為 100／300／600／1,200／6,000。
- [ ] 需要新鮮用量時，重新 fetch `/me/account/`，而不是以 task mutation 推算。

正式對外 API 細節請以
[`contracts/api-gateway/v1/public/`](../../contracts/api-gateway/v1/public/README.md)
產出的 OpenAPI 與 rules 為準。
