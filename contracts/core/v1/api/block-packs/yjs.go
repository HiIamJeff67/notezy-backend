package apicontract

import blocknote "github.com/HiIamJeff67/notezy-backend/contracts/types/blocknote"

type InitializeBlockPackYjsDocumentReqDto struct {
	Blocks []blocknote.ArborizedEditableBlock `json:"blocks"`
}

type InitializeBlockPackYjsDocumentResDto struct {
	Snapshot    []byte `json:"snapshot"`
	StateVector []byte `json:"stateVector"`
}
