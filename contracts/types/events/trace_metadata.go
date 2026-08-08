package events

type TraceMetadata struct {
	TraceParent string `json:"traceParent,omitempty"`
	TraceState  string `json:"traceState,omitempty"`
}
