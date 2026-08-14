package tokens

import "testing"

func TestGenerateAPIKeyStoresOnlyDigestMaterial(t *testing.T) {
	secret, displayPrefix, digest, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("generate API key: %v", err)
	}
	if ValidateAPIKeyFormat(secret) != nil {
		t.Fatalf("generated API key has an invalid format")
	}
	if displayPrefix == secret || len(displayPrefix) >= len(secret) {
		t.Fatalf("expected a shortened display prefix")
	}
	if digest == secret || digest != HashAPIKey(secret) {
		t.Fatalf("expected SHA-256 digest of the generated secret")
	}
}
