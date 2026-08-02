# General Go Style

## 格式與檔案

- 所有 Go 程式碼必須可被 `gofmt` 格式化；工作區已設定 Go formatter 於存檔時執行。
- 檔名採 `snake_case.go`，並依功能命名，例如 `root_shelf_service.go`、`root_shelf_controller_test.go`。
- package 名稱維持小寫單數或既有複數慣例；不要為單一用途建立新 package。
- import 由 `gofmt` 排序，標準函式庫、第三方、專案內 import 以空行分組。所有專案內 import 都必須使用顯式且準確的 package alias，例如 `dtos`、`schemas`、`exceptions`、`types`、`apitransport`、`coreadapters`；不得依賴 import path 最後一段的隱式名稱。
- 所有非空 struct literal 一律展開為多行，每個 field 各佔一行並保留 trailing comma；不可將 DTO、response 或任何只有一個 field 的 struct literal 寫在一行。這也適用於 `return` 的 struct literal 與巢狀 struct literal。

  ```go
  return &UpdateMyStationByIdResponse{
  	UpdatedAt: station.UpdatedAt,
  }, nil
  ```

## Method call 參數格式

- method/function call 的參數只在少量且簡單時保留一行。參數超過兩個、包含 nested expression、代表不同語意群組，或閱讀時需要橫向捲動/難以辨認各值角色時，必須展開為一個參數一行並保留 trailing comma。
- repository call 特別嚴格：正常 domain arguments 與 `options.With...` 是不同語意群組；有兩個以上 options 時一律多行。即使只有一個 option，只要 call 本身已有多個 domain arguments，也應多行。
- 不要把多個 `options.With...` 壓在同一行，也不要為了維持單行而省略命名清楚的 option。每個 option 佔一行，順序維持 DB/transaction、permission、soft-delete、locking、batch 等既有語意順序。

```go
station, permission, exception := s.stationRepository.CheckPermissionAndGetOneById(
	request.Body.StationId,
	request.ContextFields.UserId,
	nil,
	allowedPermissions,
	options.WithTransactionDB(tx),
	options.WithAllowedPermissions(allowedPermissions),
	options.WithOnlyDeleted(types.Ternary_Negative),
	options.WithLockingStrength(options.LockingStrengthUpdate),
)
```

## 命名與型別

- exported type/function 使用 `PascalCase`；未匯出識別字使用 `camelCase`；縮寫沿用 Go 慣例，如 `Id`、`Url`、`Db` 的專案既有寫法，不在同一領域混用兩種拼法。
- 版本化 transport DTO 使用完整 `XxxRequestDto`、`XxxResponseDto`；資料庫寫入資料使用 `XxxInput`；資料庫 table 使用 `schemas.Xxx`。`core.Request[Dto]` 與 `core.Response[Dto]` 是 Gateway 與 Core service 的 transport envelope，DTO 才是操作資料單位。
- 新功能以領域名對齊檔案、transport controller、adapter、service、repository、scope，例如 `station_*`。
- interface 僅在呼叫端需要替換實作、邊界已存在或測試確有必要時建立。新 interface 以 `XxxInterface` 對齊現有程式；不要為單一 struct 預先抽象。

## Service DTO 與 repository input 的邊界

- `internal/services/<service>/data/.../inputs` 的 `XxxInput` 是 repository persistence contract：描述 create、update、partial update 或 bulk SQL 所需的資料，只能作為 repository method 的 input。service、controller 或 gateway 不得把它當成 transport request/response contract。
- service method 若參數很多、代表同一個完整意圖，或由 Gateway 接收後要輸出到外部，就使用 `*XxxRequestDto` 作為單一 request parameter，並以 `*XxxResponseDto` 回傳資料。request / response 變數命名為 `request` / `response`；不新增不清楚的 `req`、`res` 縮寫。
- service-only 或 gateway-only workflow 也可使用具體的 `XxxRequestDto` / `XxxResponseDto` 封裝相關資料、context 與輸出，而不是用長串零散 parameters 或 anonymous struct；名稱必須描述操作，不可使用泛用的 `Data`、`Params` 或 `Payload`。
- 參數少且語意清楚時保留直接參數，不為了套用 DTO 強行包裝。只有 parameter 數量、共同 lifecycle 或呼叫邊界確實使 DTO 提高可讀性時才建立。
- service 是兩種 contract 的轉換邊界：將已驗證的 `request` 映射為 `inputs.CreateXxxInput`、`inputs.PartialUpdateXxxInput` 或 bulk input，再呼叫 repository；repository 不得 import 或依賴 transport request/response。

```go
func (s *StationService) UpdateMyStationById(
	ctx context.Context,
	request *UpdateMyStationByIdRequestDto,
) (*UpdateMyStationByIdResponseDto, *exceptions.Exception) {
	input := inputs.PartialUpdateStationInput{
		Values: inputs.UpdateStationInput{
			Name: request.Body.Values.Name,
		},
		SetNull: request.Body.SetNull,
	}

	station, exception := s.stationRepository.UpdateOneById(
		request.Body.StationId,
		request.ContextFields.UserId,
		input,
	)
	if exception != nil {
		return nil, exception
	}

	return &UpdateMyStationByIdResponseDto{
		UpdatedAt: station.UpdatedAt,
	}, nil
}
```

## 實作原則

- 先延續同一領域相鄰檔案的模式，再寫新 helper；可用既有 repository、scope、option、validator 或 exception 時不要重做。
- 函式只負責一個層級的工作。public HTTP parsing/validation 與 HTTP response 留在 Gateway controller，流程與交易留在 service，查詢與權限 SQL 留在 repository/scope。
- 使用 `context.Context` 往下傳遞：HTTP service 從 `ctx.Request.Context()` 取得，DB 以 `db.WithContext(ctx)` 建立本次操作的 session。
- 不要吞掉 error。預期的業務/基礎設施錯誤轉為 `*exceptions.Exception`，並附 `WithOrigin(err)` 保留原因。
- 優先寫直接可讀的程式，不以泛型、反射、全域狀態或新依賴解決單一功能問題。

## 批次資料庫操作（硬性規範）

**絕對不允許 per-row database operations。** `for` 只能用於整理輸入、建立 set/map、組裝 batch input、placeholder 或 response；迴圈內不得執行任何會碰觸資料庫的操作，包括：

- raw SQL：`Exec`、`Raw`、`Scan`。
- GORM query/mutation：`Model`、`Create`、`Updates`、`Update`、`Delete`、`Find`、`First`、`Count`、`Pluck` 等。
- repository method、其他 service method，或任何間接執行 SQL 的 helper。

這項規範不因資料筆數小而放寬。若需求看似只能逐筆處理，必須先改成 batch 介面、集合查詢或 SQL set-based 操作；不能在當前 function 中保留 N+1 寫法。資料庫參數或 statement 大小有限制時，可將資料切成固定大小的 batch，但每一批仍然只能是一個集合操作，不能退化成一筆一個 query。

優先順序如下：

1. 一次讀取使用 `IN ?`、join、preload、CTE 或既有 bulk repository method，將資料在記憶體中以 map/set 對應。
2. 批次新增使用 `CreateInBatches` 或既有 `CreateMany`/bulk repository method。
3. 批次更新、upsert 或關聯異動，先在迴圈中組裝 `valuePlaceholders` 與 `valueArgs`，再以單一 `VALUES`/CTE statement 執行；所有值都必須以 bound parameter 傳入，不能把使用者資料字串串進 SQL。
4. 在 SQL 需要保留輸入順序或逐筆結果時，讓 `VALUES` 帶入 index，並一次 `RETURNING`/`Scan` 結果後再於記憶體中對應。

```go
valuePlaceholders := make([]string, 0, len(bulkInputs))
valueArgs := make([]any, 0, len(bulkInputs)*2)
for _, input := range bulkInputs {
	valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::text)")
	valueArgs = append(valueArgs, input.Id, input.Name)
}

sql := fmt.Sprintf(`
	UPDATE "StationTable" AS station
	SET name = value.name
	FROM (VALUES %s) AS value(id, name)
	WHERE station.id = value.id::uuid
`, strings.Join(valuePlaceholders, ","))
result := tx.Exec(sql, valueArgs...)
```

先處理空 input，避免建立無效的 `VALUES` statement；單一 bulk query 的 `Error`、`RowsAffected` 與 transaction 收尾依既有 exception/transaction 規範處理。

## 共用函式庫

- 新增 helper 前，先檢查 [internal/shared/lib](../../internal/shared/lib/) 是否已有可直接使用的函式庫。已有相同責任的實作時必須重用，不能在 service/repository 複製一份。
- 依問題選用既有 package：去重/set 使用 `internal/shared/lib/array`，游標分頁使用 `internal/shared/lib/searchcursor`，併發工作使用 `internal/shared/lib/concurrency`，佇列與堆疊使用 `internal/shared/lib/queue`、`internal/shared/lib/stack`，以及其他已存在的 blocknote、editableblock 函式庫。跨 runtime 的 HTTP response formatting 與 public exception rendering 使用 `internal/shared/responsewriter`；rate limit 仍由各 Gateway runtime 自己持有。
- `internal/shared/lib` 僅放跨領域、可重用且與 application layer 無關的邏輯。它不可 import Notezy project code；必要的第三方 library 可以使用。僅由單一領域使用的商業規則留在該領域，不要為了「可能重用」移入 shared。
- 若既有 library 接近但不完全符合需求，優先在該 library 補最小且通用的能力；若需求只屬於單一領域，使用領域內的小 helper，避免為一次性需求建立新的 shared package。

## 變更範圍

- 一個功能變更只修改必要層與其測試；不要順手重排無關檔案或大規模改名。
- 既有未提交的使用者變更不屬於目前工作範圍，除非需求明確要求或與本次檔案有直接衝突。
- 新設定、環境變數、API 欄位、資料表欄位或事件格式都屬 contract；需要同步更新相對應文件與使用端。
