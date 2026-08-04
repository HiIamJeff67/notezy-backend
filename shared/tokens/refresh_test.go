package tokens

import "testing"

func TestRefreshTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_REFRESH_TOKEN_SECRET_KEY", "test-secret")

	token, err := GenerateRefreshToken(
		"83bdeac1-02de-42fe-a7a8-4e1a83174866",
		RefreshTokenClaims{
			Name:      "notezy",
			Email:     "notezy@example.com",
			UserAgent: "test-agent",
		},
	)
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}

	claims, err := ParseRefreshToken(*token)
	if err != nil {
		t.Fatalf("parse refresh token: %v", err)
	}
	if claims.Name != "notezy" || claims.Email != "notezy@example.com" || claims.UserAgent != "test-agent" {
		t.Fatalf("unexpected refresh token claims: %#v", claims)
	}
	if claims.Subject != "83bdeac1-02de-42fe-a7a8-4e1a83174866" {
		t.Fatalf("unexpected refresh token subject: %s", claims.Subject)
	}
	if claims.IssuedAt == nil || claims.ExpiresAt == nil {
		t.Fatal("expected generated refresh token timestamps")
	}
}
