package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const AccessTokenExpiresIn time.Duration = 30 * time.Minute

type AccessTokenClaims struct {
	Name      string `json:"name" validate:"required,min=6,max=16,alphaandnum"`
	Email     string `json:"email" validate:"required,email"`
	UserAgent string `json:"userAgent" validate:"required"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userPublicId string, claims AccessTokenClaims) (*string, error) {
	if _, err := uuid.Parse(userPublicId); err != nil {
		return nil, errors.New("access token user public ID is invalid")
	}

	secret := os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, errors.New("access token secret is required")
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userPublicId,
	}

	token, err := SignJWT(secret, claims)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func ParseAccessToken(tokenString string) (*AccessTokenClaims, error) {
	secret := os.Getenv("JWT_ACCESS_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, errors.New("access token secret is required")
	}

	claims := &AccessTokenClaims{}
	if err := ParseJWT(secret, tokenString, claims); err != nil {
		return nil, err
	}

	return claims, nil
}
