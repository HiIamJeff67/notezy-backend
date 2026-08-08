package adapterscontract

import "github.com/google/uuid"

type LoadYjsDocumentCommandDto struct{}

type LoadYjsDocumentReplyDto struct {
	Found   bool   `json:"found"`
	Payload []byte `json:"payload,omitempty"`
}

type AppendYjsUpdateCommandDto struct {
	PersistenceBatchId uuid.UUID  `json:"persistenceBatchId"`
	OriginConnectionId *uuid.UUID `json:"originConnectionId,omitempty"`
	Payload            []byte     `json:"payload"`
}

type AppendYjsUpdateReplyDto struct {
	UpdateSequence int64 `json:"updateSequence"`
}

type LoadCompactableYjsDocumentCommandDto struct{}

type LoadCompactableYjsDocumentReplyDto struct {
	Found   bool   `json:"found"`
	Payload []byte `json:"payload,omitempty"`
}

type ApplyCompactedYjsDocumentCommandDto struct {
	Payload []byte `json:"payload"`
}

type ApplyCompactedYjsDocumentReplyDto struct {
	Applied bool `json:"applied"`
}

type ApplyBlockProjectionCommandDto struct {
	Projection []byte `json:"projection"`
}

type ApplyBlockProjectionReplyDto struct {
	Applied                bool  `json:"applied"`
	ProjectedUntilSequence int64 `json:"projectedUntilSequence"`
}
