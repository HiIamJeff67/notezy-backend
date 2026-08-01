package contexts

import (
	"context"
	"testing"
)

func TestValue(t *testing.T) {
	ctx := WithValue(context.Background(), "key", "value")
	value, err := GetValue[string](ctx, "key")
	if err != nil || value != "value" {
		t.Fatalf("GetValue() = (%q, %v), want (value, nil)", value, err)
	}

	_, err = GetValue[string](ctx, "missing")
	if err == nil || err.Error() != "context value not found" {
		t.Fatalf("GetValue() error = %v, want context value not found", err)
	}
}
