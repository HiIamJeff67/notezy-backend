package blocksdto

import (
	"github.com/google/uuid"

	typescontract "github.com/HiIamJeff67/notezy-backend/contracts/types"
)

type ApplyBlockProjectionRequestDto struct {
	SchemaId          string                                 `json:"schemaId"`
	SchemaVersion     int                                    `json:"schemaVersion"`
	ProjectedSequence int64                                  `json:"projectedSequence"`
	Blocks            []typescontract.ArborizedEditableBlock `json:"blocks"`
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
