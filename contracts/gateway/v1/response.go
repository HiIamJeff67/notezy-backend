package gatewaycontract

import exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

type Response[D any] struct {
	Version   string                `json:"version"`
	Metadata  ResponseMetadata      `json:"metadata"`
	Data      D                     `json:"data"`
	Exception *exceptions.Exception `json:"exception,omitempty"`
}
