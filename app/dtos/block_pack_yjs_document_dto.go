package dtos

type InitializeBlockPackYjsDocumentReqDto struct {
	Blocks []ArborizedEditableBlock `json:"blocks"`
}

type InitializeBlockPackYjsDocumentResDto struct {
	Snapshot    []byte `json:"snapshot"`
	StateVector []byte `json:"stateVector"`
}
