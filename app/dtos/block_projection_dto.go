package dtos

type ApplyBlockProjectionInput struct {
	SchemaId          string                   `json:"schemaId"`
	SchemaVersion     int                      `json:"schemaVersion"`
	ProjectedSequence int64                    `json:"projectedSequence"`
	Blocks            []ArborizedEditableBlock `json:"blocks"`
}

type ApplyBlockProjectionResult struct {
	Applied                bool  `json:"applied"`
	ProjectedUntilSequence int64 `json:"projectedUntilSequence"`
}
