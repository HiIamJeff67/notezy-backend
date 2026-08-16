package pointers

import "testing"

func TestToPtr(t *testing.T) {
	value := "notegic"
	ptr := ToPtr(value)

	if ptr == nil {
		t.Fatal("ToPtr returned nil")
	}
	if *ptr != value {
		t.Fatalf("ToPtr() = %q, want %q", *ptr, value)
	}
}
