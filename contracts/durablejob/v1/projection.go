package durablejobcontract

import (
	"github.com/google/uuid"

	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/blocks"
)

const ApplyBlockProjectionOperation = "durablejob.apply-block-projection"

type ApplyBlockProjectionRequestDto struct {
	Documents []blocksdto.ApplyBlockProjectionDocumentRequestDto `json:"documents" validate:"required"`
}

type ApplyBlockProjectionResponseDto struct {
	AppliedBlockPackIds []uuid.UUID `json:"appliedBlockPackIds"`
}
