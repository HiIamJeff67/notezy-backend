package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

const APIKeyPrefix = "nzy_"

// GenerateAPIKey returns the one-time secret, its safe display prefix, and
// the digest that should be stored in the APIKey schema.
func GenerateAPIKey() (secret string, displayPrefix string, digest string, err error) {
	bytes := make([]byte, 32)
	if _, err = rand.Read(bytes); err != nil {
		return "", "", "", err
	}
	secret = APIKeyPrefix + base64.RawURLEncoding.EncodeToString(bytes)
	digest = HashAPIKey(secret)
	displayPrefix = secret
	if len(displayPrefix) > 12 {
		displayPrefix = displayPrefix[:12]
	}
	return secret, displayPrefix, digest, nil
}

func HashAPIKey(secret string) string {
	digest := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(digest[:])
}

func ValidateAPIKeyFormat(secret string) error {
	if !strings.HasPrefix(secret, APIKeyPrefix) || len(secret) < len(APIKeyPrefix)+32 {
		return errors.New("invalid API key format")
	}
	return nil
}
