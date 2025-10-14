package tokens

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

/* ============================== Generate Token Functions ============================== */

func generateCSRFSignature(tokenValue string) string {
	h := hmac.New(sha256.New, _csrfTokenSecret)
	h.Write([]byte(tokenValue))
	signature := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}

func GenerateCSRFToken() (*string, *exceptions.Exception) {
	randomBytes := make([]byte, _csrfTokenLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, exceptions.Token.FailedToGenerateCSRFToken().WithError(err)
	}
	tokenValue := base64.StdEncoding.EncodeToString(randomBytes)

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(_csrfTokenExpiresIn)

	signature := generateCSRFSignature(tokenValue)

	claims := types.CSRFClaims{
		Signature: signature,
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, exceptions.Token.FailedToGenerateCSRFToken().WithError(err)
	}

	token := base64.StdEncoding.EncodeToString(claimsJSON)

	return &token, nil
}

/* ============================== Parse Token Functions(only used in the validation) ============================== */

func parseCSRFToken(tokenString string) (*types.CSRFClaims, *exceptions.Exception) {
	claimsJSON, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, exceptions.Token.FailedToParseCSRFToken().WithError(err)
	}

	var claims types.CSRFClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, exceptions.Token.FailedToParseCSRFToken().WithError(err)
	}

	if time.Now().After(claims.ExpiresAt) {
		return nil, exceptions.Token.CSRFTokenExpired()
	}

	return &claims, nil
}

/* ============================== Validate Token Functions(General Purpose) ============================== */

func ValidateCSRFToken(tokenString string, expectedTokenString string) (*types.CSRFClaims, *exceptions.Exception) {
	if tokenString != expectedTokenString {
		return nil, exceptions.Token.FailedToParseCSRFToken()
	}

	claims, exception := parseCSRFToken(tokenString)
	if exception != nil {
		return nil, exception
	}

	expectedClaims, exception := parseCSRFToken(expectedTokenString)
	// this expected claim is from our backend cache,
	// so it is guarantee to be valid and we can also check if it is expired in the parse function
	if exception != nil { // will throw error if the actual CSRF token is expired
		return nil, exception
	}

	if !hmac.Equal([]byte(claims.Signature), []byte(expectedClaims.Signature)) {
		return nil, exceptions.Token.InvalidCSRFTokenSignature()
	}

	return claims, nil
}

/* ============================== Utility Functions ============================== */

func GetCSRFTokenExpiresIn() time.Duration {
	return _csrfTokenExpiresIn
}

func IsCSRFTokenExpiringSoon(claims *types.CSRFClaims) bool {
	return time.Until(claims.ExpiresAt) < time.Hour
}
