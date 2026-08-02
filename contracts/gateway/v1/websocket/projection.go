package websocketcontract

import (
	"github.com/google/uuid"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/blocks"
)

type ApplyBlockProjectionRequestDto struct {
	BlockPackId uuid.UUID                                `json:"blockPackId" validate:"required"`
	Projection  blocksdto.ApplyBlockProjectionRequestDto `json:"projection" validate:"required"`
}

type ApplyBlockProjectionResponseDto struct {
	Applied                bool  `json:"applied"`
	ProjectedUntilSequence int64 `json:"projectedUntilSequence"`
}
