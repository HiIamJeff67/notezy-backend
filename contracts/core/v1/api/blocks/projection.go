package apicontract

import (
	"github.com/google/uuid"

	blocknote "github.com/HiIamJeff67/notegic-backend/contracts/types/blocknote"
)

type ApplyBlockProjectionRequestDto struct {
	SchemaId          string                             `json:"schemaId"`
	SchemaVersion     int                                `json:"schemaVersion"`
	ProjectedSequence int64                              `json:"projectedSequence"`
	Blocks            []blocknote.ArborizedEditableBlock `json:"blocks"`
}

type ApplyBlockProjectionResponseDto struct {
	Applied                bool  `json:"applied"`
	ProjectedUntilSequence int64 `json:"projectedUntilSequence"`
}

type ApplyBlockProjectionDocumentRequestDto struct {
	BlockPackId uuid.UUID                      `json:"blockPackId"`
	Projection  ApplyBlockProjectionRequestDto `json:"projection"`
}

type ApplyBlockProjectionDocumentResponseDto []uuid.UUID
