package tokens

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestSignAndParseJWT(t *testing.T) {
	token, err := SignJWT("test-secret", jwt.MapClaims{
		"name": "notezy",
	})
	if err != nil {
		t.Fatalf("sign JWT: %v", err)
	}

	claims := jwt.MapClaims{}
	if err := ParseJWT("test-secret", token, claims); err != nil {
		t.Fatalf("parse JWT: %v", err)
	}
	if claims["name"] != "notezy" {
		t.Fatalf("unexpected JWT claims: %#v", claims)
	}
}
