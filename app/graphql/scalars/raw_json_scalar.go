package scalars

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

func MarshalRawJSON(b json.RawMessage) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if len(b) == 0 || string(b) == "null" {
			w.Write([]byte("null"))
			return
		}
		w.Write(b)
	})
}

func UnmarshalRawJSON(v interface{}) (json.RawMessage, error) {
	if v == nil {
		return nil, nil
	}

	switch v := v.(type) {
	case string:
		return json.RawMessage([]byte(v)), nil
	case []byte:
		return json.RawMessage(v), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal GraphQL JSON into json.RawMessage: %w", err)
		}
		return b, nil
	}
}
