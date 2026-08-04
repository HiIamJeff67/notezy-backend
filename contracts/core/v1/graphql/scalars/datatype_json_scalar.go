package scalars

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
	"gorm.io/datatypes"
)

func MarshalDatatypeJSON(b datatypes.JSON) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if len(b) == 0 || string(b) == "null" {
			w.Write([]byte("null"))
			return
		}
		w.Write(b)
	})
}

func UnmarshalDatatypeJSON(v interface{}) (datatypes.JSON, error) {
	if v == nil {
		return nil, nil
	}

	switch v := v.(type) {
	case string:
		return datatypes.JSON(v), nil
	case []byte:
		return datatypes.JSON(v), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal GraphQL JSON into datatypes.JSON: %w", err)
		}
		return datatypes.JSON(b), nil
	}
}
