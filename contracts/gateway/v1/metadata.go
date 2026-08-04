package gatewaycontract

import "time"

type RequestMetadata struct {
	RequestId      string `json:"requestId"`
	TraceParent    string `json:"traceParent,omitempty"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
}

type ResponseMetadata struct {
	RequestId   string    `json:"requestId"`
	RespondedAt time.Time `json:"respondedAt"`
}
