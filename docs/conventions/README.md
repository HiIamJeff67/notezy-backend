# Notezy Backend Conventions

本目錄是 Codex 與開發者修改後端時的共同基準。內容以現有程式碼為準；標示為「建議」的條目可依團隊偏好調整，調整後應直接修改本文件，讓後續代理可遵守最新版規範。

## 使用方式

1. 開始實作前，先讀與變更範圍相符的文件。
2. 既有模式與本規範衝突時，先維持既有模式；若要統一，另開明確的重構工作。
3. 新的跨模組決策應補在本目錄，而不是只留在 PR 或聊天記錄。

## 文件索引

| 文件 | 適用範圍 |
| --- | --- |
| [01-general-go.md](01-general-go.md) | Go 格式、命名、dependency direction 與變更範圍 |
| [02-architecture.md](02-architecture.md) | Gateway/API/worker ownership、HTTP request flow、GraphQL 與背景工作 |
| [03-http-api.md](03-http-api.md) | routes、controllers、adapters、request/response contract、例外與可觀測性 |
| [04-persistence.md](04-persistence.md) | services、repositories、scopes、schema、交易與 SQL |
| [05-testing-and-generated-code.md](05-testing-and-generated-code.md) | 測試、測試資料、GraphQL 產生碼與驗證清單 |
| [06-exceptions.md](06-exceptions.md) | base/service exception domain、error origin 與 `exceptions.Cover()` |

## 優先順序

1. 正確性、安全性、資料完整性與既有 API/contract。
2. 本目錄的明確規範。
3. 同一個功能附近已建立的模式。
4. Go 慣例與最小、可讀、可測的實作。

不要為「未來可能需要」新增抽象層、介面或依賴；需求出現且現有模式無法支援時再加入。

目標 workspace 與 staged migration 的 ownership 以
[microservice-architecture.md](../codebase-design/microservice-architecture.md) 為準。程式碼尚未遷移到該路徑時，不得為了符合目錄圖建立空 package 或 temporary wrapper。
