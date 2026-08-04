package blockpacksdto

import typescontract "github.com/HiIamJeff67/notezy-backend/contracts/types"

type InitializeBlockPackYjsDocumentReqDto struct {
	Blocks []typescontract.ArborizedEditableBlock `json:"blocks"`
}

type InitializeBlockPackYjsDocumentResDto struct {
	Snapshot    []byte `json:"snapshot"`
	StateVector []byte `json:"stateVector"`
}
