# Architecture and Ownership

## HTTP 功能的固定流程

```text
routes → binders → controllers → services → repositories/scopes → schemas
```

| 層 | 職責 | 不應負責 |
| --- | --- | --- |
| `routes/developmentroutes` | 註冊 URL、middleware、trace 與 metric 名稱，組裝 module | request parsing、商業流程、直接 SQL |
| `binders` | 將 header、auth context、URI/query/JSON 組成 request DTO | DB 查詢、商業決策、成功 response |
| `controllers` | 呼叫 service，統一回傳成功 JSON 或安全中止 | request parsing、權限 SQL、交易細節 |
| `services` | 驗證 DTO、協調 workflow、交易、cache/token/email 等整合 | Gin response、直接重複權限 SQL |
| `repositories` | persistence 操作、透過 scope 執行資料/權限篩選 | HTTP/Gin、回傳 HTTP status |
| `scopes` | 可重用的 GORM permission、preload、soft-delete、locking 條件 | 商業流程或 response DTO |
| `models/schemas` | table、enum、trigger、constraint、seed 的資料庫定義 | request/response 行為 |

## Module 組裝

- 每個 HTTP 領域維持 `XxxModule` 作為組裝入口，建立 scope → repository → service → binder/controller，並向 routes 暴露 `Binder` 與 `Controller`。
- 依賴由 constructor 傳入，讓單元測試能換成測試資料庫或替身；資料庫的 fallback 模式僅沿用既有 service constructor 的慣例。
- 跨領域協作由 service 透過明確依賴完成。不要讓 controller 直接呼叫另一個 controller 或 repository。

## 領域檔案與 public operation 的排序

同一領域的 public operation 有且只有一份排序。HTTP operation 以 controller interface 的順序為準；binder interface、service interface、實作 methods 與 route registration 必須逐一沿用相同順序。純 service、GraphQL 或背景工作 operation 則以其 service interface 的順序為準。新增或移動 operation 時，同一變更必須同步檢查所有相關層，不可只在其中一層隨意插入。

HTTP 操作預設依下列群組排列；同一群組內先單筆後多筆，名稱與既有 API 不同時仍遵循相同語意順序：

1. 讀取：get one、get many/search、一般 aggregate。
2. 建立：create one、create many。
3. 更新：update one、update many。
4. 還原：restore one、restore many。
5. 軟刪除：delete one、delete many。
6. 硬刪除：hard delete one、hard delete many。
7. 子資源/權限：在自己的小節中依 get、create、update/upsert、delete 排列。
8. visualization/chart：所有 `Visualize...`、chart/analytics operation 為獨立群組，不與一般 read/CRUD 混排。
9. GraphQL、system-only 或背景工作操作：與 HTTP interface 分開，置於檔案最後的明確區段。

- controller method `Xxx`、binder method `BindXxx`、service method `Xxx`、request/response DTO 與 route trace/meter 名稱必須可一眼對應；不要為同一操作創造不必要的別名。
- route registration 依同一順序排列，並以相同 module binder/controller pairing 註冊。新增 endpoint 不可只在 route 中出現而未同步其 interface 或實作。
- 有 visualization/chart operation 時，controller、binder、service interface 的 `Visualize...` methods 必須以空白行和一般 CRUD/permission methods 分組；route registration 也將 visualization endpoints 作為獨立群組排列。
- repository interface 雖不必與 HTTP operation 一對一命名，但也要依 permission/check、get、create、update、restore、soft delete、hard delete、bulk/system-only 的穩定順序分組；`System Only Method`/bulk methods 永遠放在一般 public repository API 之後。
- 舊檔案的既有順序若與本規範不同，不因無關小改動而重排。當一個操作群組被修改或新增時，維持該群組在所有涉及層的一致順序；需要全檔統一時另開明確重構工作。

每一個 controller、binder、service 或 repository 檔案，type 與 methods 的排列固定為：

```text
package / imports
interface
concrete struct（依賴欄位）
constructor（參數順序與 struct 欄位一致）
optional file-level helpers
public methods（與 interface 順序完全一致）
optional visualization/chart methods
optional GraphQL/system-only methods
```

- 每一個頂層 method declaration（包含 receiver method）之間必須有一個空白行。不可將 service、controller、repository、binder 或 route helper 的兩個 methods 緊貼在一起；空白行是 method 邊界，不以額外註解代替。
- helper 預設不建立。一段只被呼叫一次、短小且仍能在呼叫處一眼看懂的解析、mapping 或 validation 邏輯，直接在原 method 展開。不要為了少幾行而引入一次性 helper。
- 只有同一個具名概念被兩個以上 methods 重用，或單段邏輯已複雜到原地展開會遮蔽主流程時，才建立 file-level helper。例如多個 binder 都要解析相同的 station permission path parameters 時，`bindStationPermissionParams` 才是合理 helper。
- `sep30` 只在同一檔案同時存在兩個以上、可獨立導覽的 method family 時使用，例如 auxiliary helpers 與 main service methods、HTTP service methods 與 GraphQL methods、或 visualization/chart methods 與一般 CRUD methods。若只有一組 main service methods，即使它們很多，也不加 `/* ... Services for Something ... */`；順序與 methods 間的空白行已足夠。
- 若同一個 service、controller、binder、repository 或 route implementation 同時有一般 methods 與 visualization/chart methods，code implementation 必須以 `sep30` 分隔兩個 family；section label 依所在 layer 命名，例如 `/* ============================== Service Methods for Visualization ============================== */` 或 `/* ============================== Controller Methods for Visualization ============================== */`。若檔案只有 visualization methods，或只有一般 methods，則不加 `sep30`。
- 若 file-level helper 確實存在，集中放在 constructor 之後、主 methods 之前，並以 `sep30` 標示 `/* ============================== Auxiliary Functions ============================== */`。沒有 helper、GraphQL 或其他獨立 family 時，不要為唯一的 main methods 加 separator；也不要用 `sep30` 分隔每一個普通 method。
- local type、local struct、一次性 temporary variable 與匿名 function 同樣採嚴格的內嵌預設。只被引用一次、只有少量欄位、或只是為了把值轉交給下一行的資料，直接以多行 struct literal 或原地 expression 展開，不另建 type、struct 或 variable。
- local type/struct 只在無可重用的 schema/DTO/input，且它讓複雜的 query scan、bulk result mapping 或跨多個區塊的資料概念顯著清楚時使用；此時它必須在 function 內被多次使用，並以具體領域名稱命名。不能以 `Data`、`Result`、`Params` 等泛名包裝一次性資料。

### Helper 抽取為 portable lib/component

當 helper 被多個領域重用、邏輯本身夠重要且不屬於任何 application layer 時，可考慮抽成 [shared/lib](../../shared/lib/) 或獨立 component；抽取不是預設選項，必須先滿足完全獨立的條件。短小或只屬單一 layer 的 helper 仍留在原檔案原地展開。

- 新增或抽取到 `shared/lib` 的程式碼必須是 portable pure Go：只能 import Go standard library，不能 import 任何 `github.com/HiIamJeff67/notezy-backend/...` package，也不能依賴 Gin、GORM 或其他 application/framework package。
- `shared/lib` 不得接收、回傳、嵌入或 type assert application layer 型別，包括 `app/dtos`、`app/models/inputs`、schemas、repositories、services、binders、controllers、contexts、exceptions、cache client 或 `gin.Context`。它的 API 只可使用 primitive、standard library type、generic type，或由該 library 自己定義的獨立 type。
- 因此 `array.ToSet`、通用 queue/stack、與使用 generic field 的 search cursor 是合理的 shared lib；`func BuildStationResponse(req *dtos.XxxReqDto) *dtos.XxxResDto`、`func NormalizeUpdate(input inputs.PartialUpdateStationInput)` 或任何查詢 repository 的 helper 絕不可放進 shared lib。
- 若 helper 要讀寫 service DTO、repository input、DB transaction 或 HTTP/Gateway context，就代表它仍屬於該 layer，應留在 service/controller/binder/repository 中。若想抽取，先由外層把 DTO/input 轉成 neutral primitive 或 library own type，呼叫 lib 後再由外層組回 DTO/input；若無法完成這個轉換，就不要抽取。
- 獨立 component 也適用同一原則：只有其 public contract 與依賴均為 neutral 時才算獨立。若 component 必須 import DTO 或 repository，它只是 layer 內 helper，不應假裝成 reusable component。
- 可攜性的驗收標準是：把該 `shared/lib/<name>` 目錄帶到新的空白 Go module 後，只靠其目錄內程式碼與 standard library 仍可編譯/測試。既有 `shared/lib` 的 layer 依賴屬技術債，不得作為新 helper 的範例；觸及時應減少而非增加這些依賴。

- struct dependency fields、constructor parameters、constructor 初始化欄位必須維持同一順序。新增依賴時，同步更新三處，讓 wiring 可直接比對。
- module 的建立順序固定為 scope → repositories → services → binder → controller；constructor arguments 依其宣告順序傳入。`XxxModule` 的 exported fields 維持 `Binder`、`Controller` 順序，route 只透過這兩個入口連接。

## 非 HTTP 邊界

- GraphQL schema 在 `shared/graphql/schemas`，resolver 位於 `app/graphql/resolvers`。修改 schema 後依 `infra/graphql/gqlgen.yaml` 重新產生程式碼；不可手改產生檔。
- 可持久化背景工作放 `app/durablejobs/<domain>`，啟動/關閉由 `app.StartApplication` 管理。worker 要接受 `context.Context` 並支援乾淨 shutdown。
- Redis、S3、email、token、OAuth 等外部系統，使用對應的 `app/caches`、`storages`、`emails`、`tokens`、`adapters` package。商業 service 協調它們，但不複製其 transport 細節。

## 建議：新增功能前的最小檢查

1. 找同領域的 module、route、binder、DTO、service、repository 和 test。
2. 判斷是否為既有功能擴充；若是，優先在同一條 flow 增加一個操作。
3. 只有新資源、獨立 lifecycle 或明確外部 contract 才建立完整新 module。
