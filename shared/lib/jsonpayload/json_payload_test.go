package jsonpayload

import "testing"

func TestDecode(t *testing.T) {
	var payload struct {
		Name string `json:"name"`
	}
	if err := Decode([]byte(`{"name":"notegic"}`), &payload); err != nil {
		t.Fatalf("decode JSON payload: %v", err)
	}
	if payload.Name != "notegic" {
		t.Fatalf("unexpected decoded payload: %#v", payload)
	}
	if err := Decode([]byte(" \t"), &payload); err == nil {
		t.Fatal("expected empty payload error")
	}
}
