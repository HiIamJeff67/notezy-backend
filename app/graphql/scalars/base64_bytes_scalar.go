package scalars

import (
	"encoding/base64"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalBase64Bytes(b []byte) graphql.Marshaler {
	return graphql.MarshalString(base64.StdEncoding.EncodeToString(b))
}

func UnmarshalBase64Bytes(v interface{}) ([]byte, error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("invalid bytes: %v, should be base64 string", v)
	}
	return base64.StdEncoding.DecodeString(s)
}
