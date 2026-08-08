package gatewaycontract

import "encoding/json"

/* ============================== Request ============================== */

type Request[D any] struct {
	Version   string          `json:"version"`
	Operation string          `json:"operation"`
	Metadata  RequestMetadata `json:"metadata"`
	Tokens    Tokens          `json:"tokens,omitempty"`
	Dto       D               `json:"dto"`
}

// Tokens carries authentication credentials across the Gateway/Core boundary.
// Gateway owns cookie extraction; Core only receives this typed value.
type Tokens struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	CSRFToken    string `json:"csrfToken,omitempty"`
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

/* ============================== Client Request ============================== */

// ClientRequest is the public request envelope exposed by the Gateway.
// Internal runtimes must use Request instead.
type ClientRequest[D any] struct {
	Dto D `json:"dto"`
}

func (r ClientRequest[D]) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Dto)
}

func (r *ClientRequest[D]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Dto)
}
