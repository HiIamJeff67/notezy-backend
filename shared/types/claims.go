package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Id        string `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required,min=6,max=16,alphaandnum"`
	Email     string `json:"email" validate:"required,email"`
	UserAgent string `json:"userAgent" validate:"required"`
	jwt.RegisteredClaims
}

type CSRFClaims struct {
	Signature string    `json:"signature" validate:"required"`
	ExpiresAt time.Time `json:"expiresAt"`
	IssuedAt  time.Time `json:"issuedAt"`
}
