package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"time"
)

const CSRFTokenExpiresIn time.Duration = 7 * 24 * time.Hour

type CSRFTokenClaims struct {
	Signature string    `json:"signature" validate:"required"`
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt  time.Time `json:"issuedAt"`
}

func GenerateCSRFToken(claims CSRFTokenClaims) (*string, error) {
	secret := os.Getenv("CSRF_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, errors.New("csrf token secret is required")
	}

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err
	}
	tokenValue := base64.StdEncoding.EncodeToString(randomBytes)
	issuedAt := time.Now()
	claims.Signature = signCSRFToken(secret, tokenValue)
	claims.ExpiresAt = issuedAt.Add(CSRFTokenExpiresIn)
	claims.IssuedAt = issuedAt
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	token := base64.StdEncoding.EncodeToString(claimsJSON)
	return &token, nil
}

func ValidateCSRFToken(tokenString string, expectedTokenString string) (*CSRFTokenClaims, error) {
	secret := os.Getenv("CSRF_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, errors.New("csrf token secret is required")
	}
	if tokenString != expectedTokenString {
		return nil, errors.New("csrf token is inconsistent")
	}

	claims, err := parseCSRFToken(tokenString)
	if err != nil {
		return nil, err
	}
	expectedClaims, err := parseCSRFToken(expectedTokenString)
	if err != nil {
		return nil, err
	}
	if !hmac.Equal([]byte(claims.Signature), []byte(expectedClaims.Signature)) {
		return nil, errors.New("csrf token signature is invalid")
	}

	return claims, nil
}

func IsCSRFTokenExpiringSoon(claims *CSRFTokenClaims) bool {
	return time.Until(claims.ExpiresAt) < time.Hour
}

func signCSRFToken(secret string, tokenValue string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	_, _ = hash.Write([]byte(tokenValue))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func parseCSRFToken(tokenString string) (*CSRFTokenClaims, error) {
	claimsJSON, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, err
	}

	var claims CSRFTokenClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, err
	}
	if time.Now().After(claims.ExpiresAt) {
		return nil, errors.New("csrf token is expired")
	}

	return &claims, nil
}
