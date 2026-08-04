package coreapicontract

// RequestDto standardizes the sections carried by a client request.
type RequestDto[Header any, Body any, Param any, Query any] struct {
	Header Header `json:"header"`
	Body   Body   `json:"body"`
	Param  Param  `json:"param"`
	Query  Query  `json:"query"`
}
