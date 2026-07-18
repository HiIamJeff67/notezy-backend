package dtos

import "github.com/google/uuid"

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

type ApplyBlockProjectionDocumentInput struct {
	BlockPackId uuid.UUID
	Projection  ApplyBlockProjectionInput
}

type ApplyBlockProjectionDocumentResult []uuid.UUID
