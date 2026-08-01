package tokens

import "testing"

func TestCSRFTokenRoundTrip(t *testing.T) {
	t.Setenv("CSRF_TOKEN_SECRET_KEY", "test-secret")

	token, err := GenerateCSRFToken(
		CSRFTokenClaims{},
	)
	if err != nil {
		t.Fatalf("generate csrf token: %v", err)
	}

	claims, err := ValidateCSRFToken(*token, *token)
	if err != nil {
		t.Fatalf("validate csrf token: %v", err)
	}
	if claims.Signature == "" || claims.IssuedAt.IsZero() || claims.ExpiresAt.IsZero() {
		t.Fatalf("unexpected csrf token claims: %#v", claims)
	}
}
