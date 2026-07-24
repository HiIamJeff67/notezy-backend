package dtos

import "github.com/google/uuid"

type ApplyBlockProjectionReqDto struct {
	SchemaId          string                   `json:"schemaId"`
	SchemaVersion     int                      `json:"schemaVersion"`
	ProjectedSequence int64                    `json:"projectedSequence"`
	Blocks            []ArborizedEditableBlock `json:"blocks"`
}

type ApplyBlockProjectionResDto struct {
	Applied                bool  `json:"applied"`
	ProjectedUntilSequence int64 `json:"projectedUntilSequence"`
}

type ApplyBlockProjectionDocumentReqDto struct {
	BlockPackId uuid.UUID
	Projection  ApplyBlockProjectionReqDto
}

type ApplyBlockProjectionDocumentResDto []uuid.UUID
