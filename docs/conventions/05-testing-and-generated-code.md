# Tests, Generated Code, and Verification

## 測試位置與寫法

- 與 package 緊密相關的單元測試放在同 package，檔名為 `*_test.go`；跨 HTTP 邊界的流程測試放 `test/e2e/<domain>`。
- 新增或修正 binder 時，使用 `httptest` + Gin test router 驗證 request parsing、context 欄位與不合法輸入；沿用 `app/binders/*_test.go` 的模式。
- service/repository 的測試先覆蓋權限、soft delete、交易和錯誤分支等風險高的行為，而非為每個 getter 建立樣板測試。
- 可重用的 JSON fixture 放在同領域 `testdata/`，以描述行為的 `snake_case.json` 命名。
- 測試須獨立、可重複執行，不依賴執行順序、實際時鐘、外部網路或共享 production 資料。

## Generated code 與契約

- `app/graphql/generated` 與 `app/graphql/models/models_gen.go` 是 GraphQL 產生物；修改 GraphQL schema 後以 `make gql-generate` 或 `make gql-regenerate` 更新，絕不手動編輯生成碼。
- API、GraphQL、realtime protocol、資料庫行為等跨端約定放 `docs/contracts/`。改變這些公開語意時，同一變更需更新 contract 和相關測試。
- `infra/` 是部署/監控設定，變更後檢查 Docker Compose、Nginx、OTEL/Grafana 設定是否仍彼此一致。

## 每次變更的最小驗證

1. 執行 `gofmt` 於修改的 Go 檔案。
2. 執行受影響 package 的 `go test`；若碰到可安全啟動的 integration/e2e 環境，再執行相對應 e2e test。
3. 若改 schema，產生 GraphQL code 或 DB migration，確認產物與註冊檔都已更新。
4. 若改 route、DTO 或 exception，檢查前端/contract 的 payload、HTTP status 與 metric 名稱。

當環境或既有問題使測試無法執行，提交說明必須列出未執行的指令、原因與可能影響。
