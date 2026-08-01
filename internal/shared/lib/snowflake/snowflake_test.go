package snowflake

import "testing"

func TestGenerateUniqueID(t *testing.T) {
	firstID := GenerateUniqueID()
	secondID := GenerateUniqueID()
	if firstID == secondID {
		t.Fatal("expected unique IDs")
	}
}
