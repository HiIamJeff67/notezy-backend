# Exception Conventions

## 基本原則

- 所有可預期的 application、validation、permission、persistence 與外部依賴失敗，使用 `app/exceptions` 的既有 domain factory；不要在 request flow 中自行建立另一種 JSON error 格式。
- 選擇最接近資源/操作的 domain，例如 shelf、station、routine；底層 Go error 以 `WithOrigin(err)` 附加，保留 log 與 trace 所需的實際原因。
- binder 負責將 bind/parse 失敗轉為 `InvalidInput()` 或 `InvalidDto()`；service 驗證 request DTO 時使用其領域的 `InvalidDto()`；repository 將 GORM/raw SQL 結果轉為該資料領域的 `NotFound()`、`FailedToCreate()`、`FailedToUpdate()`、`NoChanges()` 等 exception。
- exception 往上回傳，不在 service/repository 私自轉成 HTTP response。只有 controller、binder 或 middleware 使用安全 JSON response 方法；internal exception 不能直接暴露給 client。
- 產生 exception 與緊接的 return/rollback 是同一錯誤處理區塊，保持相鄰，除非中間確實有必要的補償動作。

## `exceptions.Cover()`

`exceptions.Cover(existing, fallbacks)` 只在 `existing == nil` 時依序取第一個成立的 fallback；已有 exception 時原樣回傳。因此它適合把「同一個操作可能失敗的多個條件」集中為一個 guard，減少重複的 `if`、rollback 和 return。

- 兩個以上彼此替代的結果條件時使用，例如 GORM `Error`、查無資料、`RowsAffected == 0`，或 repository exception 之後的領域不變量驗證。
- fallback 順序就是優先順序。先列出帶有 `WithOrigin(result.Error)` 的底層錯誤，再列空結果、無變更或領域條件，避免丟失較具體的原因。
- 只有一個清楚條件時直接寫普通 `if`；不要為單一 guard 或無關的錯誤流程強行使用 `Cover()`。
- `Cover()` 只挑選 exception，不會執行 rollback、log 或 response；在取得非 nil exception 後，立刻走該 scope 的既有錯誤收尾。

```go
result := parsedOptions.DB.Model(&schemas.Station{}).
	Where("id = ? AND deleted_at IS NULL", id).
	Select("*").
	Updates(&updates)
if exception := exceptions.Cover(nil, []types.Pair[bool, *exceptions.Exception]{
	{First: result.Error != nil, Second: exceptions.Station.FailedToUpdate().WithOrigin(result.Error)},
	{First: result.RowsAffected == 0, Second: exceptions.Station.NoChanges()},
}); exception != nil {
	parsedOptions.DB.Rollback()
	return nil, exception
}
```

當前面呼叫已經回傳 exception，而後面還要補資源關係、ownership 或上限等不變量時，將既有 exception 作為第一個參數。若 repository 已失敗，fallback 不會覆寫它；若 repository 成功，才檢查後續條件。

```go
station, exception := r.CheckPermissionAndGetOneById(id, userId, opts...)
if exception = exceptions.Cover(exception, []types.Pair[bool, *exceptions.Exception]{
	{First: station.OwnerId != expectedOwnerId, Second: exceptions.Station.NotFound()},
	{First: station.ArchivedAt != nil, Second: exceptions.Station.NotFound()},
}); exception != nil {
	parsedOptions.DB.Rollback()
	return nil, exception
}
```

`Cover()` 的條件應是安全可求值的：若前一步可能回傳 nil resource，先讓 repository exception 在前面被處理，或確保 fallback 不會解參考 nil pointer。
