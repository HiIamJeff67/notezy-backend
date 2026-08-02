package websocketcontract

import "github.com/google/uuid"

type LoadCompactableYjsDocumentRequestDto struct {
	BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
}

type LoadCompactableYjsDocumentResponseDto struct {
	Found   bool   `json:"found"`
	Payload []byte `json:"payload"`
}

type ApplyCompactedYjsDocumentRequestDto struct {
	BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
	Payload     []byte    `json:"payload" validate:"required"`
}

type ApplyCompactedYjsDocumentResponseDto struct {
	Applied bool `json:"applied"`
}

type LoadYjsDocumentRequestDto struct {
	BlockPackId uuid.UUID `json:"blockPackId" validate:"required"`
}

type LoadYjsDocumentResponseDto struct {
	Payload []byte `json:"payload"`
}

type AppendYjsUpdateRequestDto struct {
	BlockPackId        uuid.UUID  `json:"blockPackId" validate:"required"`
	PersistenceBatchId uuid.UUID  `json:"persistenceBatchId" validate:"required"`
	OriginConnectionId *uuid.UUID `json:"originConnectionId"`
	Payload            []byte     `json:"payload" validate:"required"`
}

type AppendYjsUpdateResponseDto struct {
	UpdateSequence int64 `json:"updateSequence"`
}
