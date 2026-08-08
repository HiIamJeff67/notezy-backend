# Tests, Generated Code, and Verification

## 測試位置與寫法

- 與 package 緊密相關的單元測試放在 production 檔案旁的同 package，檔名為 `*_test.go`；runtime 對外或對另一 runtime 的 e2e 測試放在該 runtime 的 `test/e2e/`。例如 status endpoint 的 HTTP 行為測試放在 `internal/<runtime>/test/e2e/status/status_test.go`，而不是 `transports/status/`；跨 runtime 的 integration 與 architecture 測試放在獨立的 `test/` module 內。
- `test/e2e/` 只驗證跨 package、HTTP/WebSocket 邊界或 runtime 間的完整流程；它不承擔一般 unit test。單一 utility、parser 或 service 的純邏輯仍與 production package 同層測試；需要透過 HTTP、Gin、ServeMux 或 WebSocket 驗證 exposed endpoint 的測試則歸入 runtime 的 e2e。
- `contracts/`、`shared/` 與各 runtime 只保留自己的 unit tests；`test/go.mod` 管理跨 module 測試依賴與本地 `replace`，不得讓這些測試依賴 root module。
- 新增或修正 Gateway controller 時，使用 `httptest` + Gin test router 驗證 request parsing、actor/context 欄位與不合法輸入。
- service/repository 的測試先覆蓋權限、soft delete、交易和錯誤分支等風險高的行為，而非為每個 getter 建立樣板測試。
- 可重用的 JSON fixture 放在同領域 `testdata/`，以描述行為的 `snake_case.json` 命名。
- 測試須獨立、可重複執行，不依賴執行順序、實際時鐘、外部網路或共享 production 資料。

## Generated code 與契約

- `contracts/core/v1/graphql/generated` 與 `contracts/core/v1/graphql/models` 是 GraphQL 產生物；修改 `contracts/core/v1/graphql` source 後以 `make -C contracts gql-generate` 或 `make -C contracts gql-regenerate` 更新，絕不手動編輯生成碼。gqlgen config 的 generated output 必須維持在這兩個 contracts path。
- API route 公開語意放 `docs/api-route-design/`；程式碼與資料模型設計放 `docs/codebase-design/`；Realtime、Yjs 與跨 runtime 協定放 `docs/system-design/`。改變任何公開語意時，同一變更必須更新對應設計文件與相關測試。
- `infra/` 是部署/監控設定，變更後檢查 Docker Compose、Nginx、OTEL/Grafana 設定是否仍彼此一致。

## Go module 與測試入口

- repository root 只保留 `go.work` 作為 workspace manifest，不保留 root `go.mod`／`go.sum`。
- `make test-all` 會透過 `internal/cli` 依序執行 contracts、shared、各 Go runtime 與 root `test/` module；要測試單一 module，使用 `make test-module MODULE=<module>` 或直接在該 module 目錄執行 `GOWORK=off go test ./...`。
- `make test-architecture` 與 `make test-integration` 只執行 root `test/` module 的對應分層；runtime E2E 由所屬 runtime 的 Makefile 執行。
- Root integration tests 使用 `testcontainers-go` 建立一次性 PostgreSQL、Redis 與 Kafka；預設只編譯並跳過 Docker-backed cases，設定 `NOTEZY_RUN_INTEGRATION=1` 才會啟動容器。測試資料放在 `test/integration/testdata/`。Core／DurableJob Kafka broker flow 使用 `make test-integration-kafka` 執行；可用 `KAFKA_BROKERS` 指定既有 broker，未指定時由測試啟動一次性 Kafka Testcontainer。
- WebSocket load／soak 與 Kafka consumer-lag verification 使用 `test/load/` 下的 Grafana k6 scripts，透過 `make test-load-websocket`、`make test-soak-websocket` 與 `make test-load-kafka-lag` 執行；這些腳本不進入任何 Go module。
- GraphQL 生成命令由 `contracts/Makefile` 維護，使用 `make -C contracts gql-generate` 從 contracts module 執行；產物仍寫入 `contracts/core/v1/graphql/`。不得因工具命令而恢復 root `go.mod`。
- root Makefile 只負責委派；module-specific 命令放在 `contracts/`、`shared/`、各 runtime 與 `test/` 的 Makefile。`internal/cli` 使用 Cobra 與 subprocess 呼叫這些 module 命令，不直接 import 測試 package。

## 每次變更的最小驗證

1. 執行 `gofmt` 於修改的 Go 檔案。
2. 執行受影響 package 的 `go test`；若碰到可安全啟動的 integration/e2e 環境，再執行相對應 e2e test。
3. 若改 schema，產生 GraphQL code 或 DB migration，確認產物與註冊檔都已更新。
4. 若改 route、DTO 或 exception，檢查前端/contract 的 payload、HTTP status 與 metric 名稱。

當環境或既有問題使測試無法執行，提交說明必須列出未執行的指令、原因與可能影響。
