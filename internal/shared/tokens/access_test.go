package tokens

import "testing"

func TestAccessTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_ACCESS_TOKEN_SECRET_KEY", "test-secret")

	token, err := GenerateAccessToken(
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		AccessTokenClaims{
			Name:      "notezy",
			Email:     "notezy@example.com",
			UserAgent: "test-agent",
		},
	)
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	claims, err := ParseAccessToken(*token)
	if err != nil {
		t.Fatalf("parse access token: %v", err)
	}
	if claims.Name != "notezy" || claims.Email != "notezy@example.com" || claims.UserAgent != "test-agent" {
		t.Fatalf("unexpected access token claims: %#v", claims)
	}
	if claims.Subject != "83bdeac1-02de-42fe-a7a8-4e1a83174866" {
		t.Fatalf("unexpected access token subject: %s", claims.Subject)
	}
	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		t.Fatal("expected generated access token timestamps")
	}
}
