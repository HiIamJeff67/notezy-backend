package blockpacksdto

import (
	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/api/v1/blocks"
)

type InitializeBlockPackYjsDocumentReqDto struct {
	Blocks []blocksdto.ArborizedEditableBlock `json:"blocks"`
}

type InitializeBlockPackYjsDocumentResDto struct {
	Snapshot    []byte `json:"snapshot"`
	StateVector []byte `json:"stateVector"`
}
