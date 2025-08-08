package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id        string `json:"id"`
	Name      string `json:"name" validate:"required,min=6,max=16,alphaandnum"`
	Email     string `json:"email" validate:"required,email"`
	UserAgent string `json:"userAgent" validate:"required"`
	jwt.RegisteredClaims
}
