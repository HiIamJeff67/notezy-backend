package gatewaycontract

type Request[D any] struct {
	Version   string          `json:"version"`
	Operation string          `json:"operation"`
	Metadata  RequestMetadata `json:"metadata"`
	Dto       D               `json:"dto"`
}

func (r *Request[D]) GetVersion() string {
	return r.Version
}

func (r *Request[D]) GetOperation() string {
	return r.Operation
}

func (r *Request[D]) GetMetadata() RequestMetadata {
	return r.Metadata
}
