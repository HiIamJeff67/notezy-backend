package adapterscontract

type BlockPackQuotaPolicy struct {
	Version           int   `json:"version"`
	MaximumBlockCount int32 `json:"maximumBlockCount"`
}
