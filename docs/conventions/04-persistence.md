# Service and Persistence Conventions

## Service

- service 是商業 workflow 與 transaction boundary。每個公開 HTTP workflow 先用 `validation.Validator.Struct(reqDto)` 驗證，再用 `s.db.WithContext(ctx)` 取得 request-scoped DB。
- service 回傳 response DTO 或必要領域資料與 `*exceptions.Exception`，不回傳 `gin.Context`、HTTP status 或 `gin.H`。
- cache、email、token、storage 和 realtime 操作應在 service 的 workflow 中明確排序。失敗補償與可重試語意要與資料庫提交關係一致。

## Service method 的區塊與留白

一個 service method 應依「同一個語意區塊緊鄰，不同階段以一個空白行分隔」閱讀，而不是依行數或固定模板切割。常見順序是：輸入驗證 → 前置條件/準備資料 → 建立 DB session 或 transaction → 讀寫 workflow → commit → response。沒有發生的階段不需要補空白或註解。

- 一次操作與緊接的錯誤處理屬於同一語意區塊：repository/DB 呼叫後直接接 `if exception != nil` 或 `if err != nil`，中間不留空行。
- validation guard 位於 method 最上方；guard 結束後留一個空白行，再開始下一個獨立階段。只有在 validation 後立即接同一個語意的前置條件時，才可不分隔。
- `db := s.db.WithContext(ctx)` 或 `tx := s.db.WithContext(ctx).Begin()` 是 workflow 的執行環境，應自成一個區塊：與前一個階段及其後第一個 query/workflow 各以空白行分開。
- 衍生值、input 組裝、查詢結果轉 DTO、迴圈處理等，各自形成連續區塊；切換到另一種工作時以一個空白行分隔。不要在同一區塊內為了視覺效果插入空白行。
- 一般 method 以空白行表達階段即可。`sep30` 只在同一檔案有兩個以上明確且複雜的獨立 method family（例如 helper、HTTP service、chart service、GraphQL service）時才可使用；若檔案只有 main service methods，不得在它們上方加 `Services for Something` separator：

  ```go
  /* ============================== Service Methods for GraphQL Station ============================== */
  ```

  `sep30` 用於多個 file-level family 的導覽，不用來分隔單一 method 的每一步、每個 `if` 或每個 struct field，也不用來為唯一的 main methods 加標題；一個同類型的 service 檔案不需要它。

## Service validation 與 DB query 格式

每個接收 request DTO 的 service method 一開始先驗證 DTO，並將 validator error 轉成該領域的 exception。接著才建立 request-scoped DB；兩者是不同階段，所以維持上下空白行。

```go
func (s *StationService) CreateStation(
	ctx context.Context,
	reqDto *dtos.CreateStationReqDto,
) (*dtos.CreateStationResDto, *exceptions.Exception) {
	if err := validation.Validator.Struct(reqDto); err != nil {
		return nil, exceptions.Station.InvalidDto().WithOrigin(err)
	}

	db := s.db.WithContext(ctx)

	newStationId, exception := s.stationRepository.CreateOne(
		reqDto.ContextFields.UserId,
		input,
		options.WithDB(db),
	)
	if exception != nil {
		return nil, exception
	}

	return &dtos.CreateStationResDto{
		Id: *newStationId, 
	}, nil
}
```

GORM chain 超過一個操作時換行書寫；receiver 的 `.` 留在行尾，讓每個 `Model`、`Where`、`Order`、`Find` 或 `Updates` 都是清楚的一步。query 建立與其執行屬於同一區塊；result/error 判斷緊接在後，不在中間插空白行。

```go
var blocks []schemas.Block
if err := db.Model(&schemas.Block{}).
	Where("block_pack_id = ?", reqDto.Param.BlockPackId).
	Order("created_at ASC").
	Order("id ASC").
	Find(&blocks).Error; err != nil {
	return nil, exceptions.Block.NotFound().WithOrigin(err)
}

result := tx.
	Model(&schemas.Station{}).
	Where("id = ?", station.Id).
	Update("deleted_at", time.Now())
if result.Error != nil {
	tx.Rollback()
	return nil, exceptions.Station.FailedToUpdate().WithOrigin(result.Error)
}
```

短且單一的 DB 操作可留在一行；不要為了套用鏈式格式而把簡單的 `tx.Model(&schema).Update(...)` 刻意拆開。

## Service transaction 與 repository transaction

多個讀寫必須一起成功或一起失敗時，transaction 由 service 開啟並由 service 唯一負責 `Commit` / `Rollback`。傳給 repository 的每一次呼叫都使用同一個 `tx` 和 `options.WithTransactionDB(tx)`；該 option 會同時帶入 DB 並標記 transaction 已開始，避免 repository 因 `IsTransactionStarted` 為 false 而另開 nested transaction。

`tx` 的建立獨立成一個區塊。若 `Begin()` 回傳 error，尚未有可回滾的交易，直接回傳；開始後的任何失敗，rollback 與 return 是同一個錯誤收尾區塊，兩者必須相鄰、不插空白行或其他程式碼。

```go
tx := s.db.WithContext(ctx).Begin()
if err := tx.Error; err != nil {
	return nil, exceptions.Shelf.FailedToCommitTransaction("failed to begin transaction").WithOrigin(err)
}

rootShelf, exception := s.rootShelfRepository.CheckPermissionAndGetOneById(
	reqDto.Body.RootShelfId,
	reqDto.ContextFields.UserId,
	nil,
	allowedPermissions,
	options.WithTransactionDB(tx),
	options.WithOnlyDeleted(types.Ternary_Negative),
	options.WithLockingStrength(options.LockingStrengthUpdate),
)
if exception != nil {
	tx.Rollback()
	return nil, exception
}

// All remaining reads and writes use tx or options.WithTransactionDB(tx).

if err := tx.Commit().Error; err != nil {
	tx.Rollback()
	return nil, exceptions.Shelf.FailedToCommitTransaction().WithOrigin(err)
}
```

當單一 repository method 本身擁有完整的 atomic workflow 時，沿用其既有的「repository 自行開/提交 transaction」行為；service 不可在外層 transaction 中混用另一個 `s.db` 或另一個 transaction。新跨 repository workflow 一律由 service 管理外層 transaction。

## Repository 與 scope

- repository 集中 GORM/raw SQL 的存取，公開方法以動作清楚命名，例如 `GetOneById`、`CreateMany`、`UpdateOneById`。
- query 請使用 `schemas.Xxx` model 與既有 `scope` 封裝 permission、preload、soft-delete 和 locking；不要在 service/controller 重複手寫存取控制的 `Where` 條件。
- repository option 透過 `options.WithDB`、`WithOnlyDeleted`、`WithLockingStrength` 等既有 option 傳入。只有確實的維護/內部需求才可使用 `WithSkipPermissionCheck()`，且呼叫處必須能說明為何安全。
- 輸入資料使用 `models/inputs` 的 create/update/partial-update 型別。不要直接把 request DTO 餵給 GORM。
- GORM result error 要轉為對應領域 exception 並保留 origin；`First`、`Find` 後也要處理空結果的領域語意，不能只看 `result.Error`。

## Repository partial update

需要支援只更新部分欄位或明確設為 `NULL` 時，使用既有的 partial update flow；不要自行以 map 拼接欄位，也不要把 request DTO 直接交給 `Updates`。

1. 在該領域的 `app/models/inputs/<domain>_input.go` 定義 `UpdateXxxInput`，欄位使用 pointer 表示「本次有提供值」，保留正確的 `json` 與 `gorm:"column:..."` tag。
2. 將 input 接到 [partial_update_input.go](../../app/models/inputs/partial_update_input.go)：`type PartialUpdateXxxInput = PartialUpdateInput[UpdateXxxInput]`。`Values` 載有要覆寫的值，`SetNull` 表示需要設為 `NULL` 的欄位。
3. repository 在同一筆 transaction 中先取得已存在且已通過 permission check 的 schema，完成任何關聯資源/ownership 驗證後，呼叫 `util.PartialUpdatePreprocess(input.Values, input.SetNull, *existing)`。
4. 將合併結果用 `Select("*").Updates(&updates)` 寫回；因 processor 已把未提供的欄位保留為既有值，`Select("*")` 才能正確寫入明確要求的零值或 `NULL`。

```go
type UpdateStationInput struct {
	Name        *string `json:"name" gorm:"column:name;"`
	Description *string `json:"description" gorm:"column:description;"`
}

type PartialUpdateStationInput = PartialUpdateInput[UpdateStationInput]
```

```go
updates, err := util.PartialUpdatePreprocess(input.Values, input.SetNull, *existingStation)
if err != nil {
	parsedOptions.DB.Rollback()
	return nil, exceptions.Util.FailedToPreprocessPartialUpdate(
		input.Values,
		input.SetNull,
		*existingStation,
	).WithOrigin(err)
}

result := parsedOptions.DB.Model(&schemas.Station{}).
	Where("id = ? AND deleted_at IS NULL", id).
	Select("*").
	Updates(&updates)
```

- `SetNull` 的 key 使用對應 Go field name；processor 會處理大小寫與底線差異，但新 API/DTO 應仍維持既有的 camelCase contract。
- partial update 前的關聯驗證只在該欄位有新值且未標記為 `SetNull` 時執行，例如移動 parent resource 前先檢查目的地 permission。
- processor error、DB error 與 `RowsAffected == 0` 的結果，使用 [exception conventions](06-exceptions.md) 的 domain exception 與 `exceptions.Cover()` 處理；在自己擁有 transaction 的 repository 中，rollback 和 return 緊鄰。

## Schema、migration 與 SQL

- table schema 放 `app/models/schemas`；enum、constraint、trigger、seed、raw SQL 放各自既有子目錄。
- 新 table/enum/trigger/constraint 必須註冊至其對應的 `migrate.go`，否則 migration 不會套用。
- soft-delete、ownership、projection/accounting 等資料庫不變量，優先延續現有 trigger/constraint/scope 模式；不可只依賴 controller 的檢查。
- 修改 trigger 或 raw SQL 時，確認引用到的表名、欄位與 migration 註冊都同步；資料庫契約改動也要更新相關 `docs/contracts`。
