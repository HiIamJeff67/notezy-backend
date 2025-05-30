package global

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
