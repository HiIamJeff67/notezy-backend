package jsonpayload

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Decode(payload []byte, target any) error {
	if len(bytes.TrimSpace(payload)) == 0 {
		return fmt.Errorf("payload is required")
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return fmt.Errorf("decode JSON payload: %w", err)
	}

	return nil
}
