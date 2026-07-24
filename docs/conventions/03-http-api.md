# HTTP API Conventions

## Routes 與 middleware

- 主 API route 放在 `app/routes/developmentroutes`；僅測試需要的 route 放在 `app/routes/testroutes`。
- route 以該領域的 `NewXxxModule()` 組裝，並將 binder 包住 controller method。
- 延續該 route 既有的 middleware 順序與 `RepositionMiddleware` 用法。每個新增 endpoint 必須有對應的 trace operation 名稱與 `server.requests.<domain>.<operation>` metric 名稱。
- 權限由 middleware 與 repository scope 共同保護。route 使用 `AllowedPermissionsAbove` 設定入口門檻；資料讀寫仍必須透過帶 permission check 的 repository/scope，避免繞過保護。
- 維持既有 URL 與 payload 命名的 camelCase。新 API 的 method、path 與 response 形狀需要與前端 contract 一起決定；不要在無相容方案下改舊欄位。

## DTO 與 binder

- request DTO 定義在 `app/dtos`，嵌入 `NotezyRequest[Header, ContextFields, Body, Param]`；只使用該 endpoint 需要的四個部分，沒有資料時用 `any`。
- JSON 使用 `json` tag，query 使用 `form` tag，path 參數使用 `uri` tag，並以 `validate` tag 表達輸入限制。
- binder 負責 `ShouldBindJSON`、`ShouldBindQuery`、`ctx.Param`、`ctx.Query` 與 auth context 轉型；轉型/綁定失敗要回傳此領域的 `InvalidInput()` 或 `InvalidDto()` exception。
- DTO 只能描述 transport contract；不要把 GORM model、DB handle、Gin context 或商業邏輯放入 DTO。

## Controller 與 response

- controller 只取得 `ctx.Request.Context()`、呼叫 service、處理 `*exceptions.Exception`，最後回傳成功 response。
- 成功格式保持 `{ "success": true, "data": ..., "exception": null }`；失敗一律使用 exception 的安全 response 方法，保持 `{ "success": false, "data": null, "exception": ... }`。
- 不要在 controller 中判斷 ownership、組 GORM query、做交易或直接讀取 request body。

## Exception、log 與觀測性

- 使用 `app/exceptions` 的領域 exception factory；不要自行建立臨時 HTTP error 格式或直接暴露 internal error。
- DB、cache、parser 等底層 error 以 `WithOrigin(err)` 保留，並用 `Log()` / 安全 response 走既有統一流程。internal exception 不得直接回傳給客戶端。
- 新 endpoint 對應 trace/metric；重大背景流程也應沿用 `app/monitor` 的 logger、meter、tracer，而非另建 logging framework。
