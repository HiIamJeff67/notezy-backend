package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const RefreshTokenExpiresIn time.Duration = 7 * 24 * time.Hour

type RefreshTokenClaims struct {
	Name      string `json:"name" validate:"required,min=6,max=16,alphaandnum"`
	Email     string `json:"email" validate:"required,email"`
	UserAgent string `json:"userAgent" validate:"required"`
	jwt.RegisteredClaims
}

func GenerateRefreshToken(userPublicId string, claims RefreshTokenClaims) (*string, error) {
	if _, err := uuid.Parse(userPublicId); err != nil {
		return nil, errors.New("refresh token user public ID is invalid")
	}

	secret := os.Getenv("JWT_REFRESH_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, errors.New("refresh token secret is required")
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Subject:   userPublicId,
	}

	token, err := SignJWT(secret, claims)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func ParseRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	secret := os.Getenv("JWT_REFRESH_TOKEN_SECRET_KEY")
	if secret == "" {
		return nil, errors.New("refresh token secret is required")
	}

	claims := &RefreshTokenClaims{}
	if err := ParseJWT(secret, tokenString, claims); err != nil {
		return nil, err
	}

	return claims, nil
}
