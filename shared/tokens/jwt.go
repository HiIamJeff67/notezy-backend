package tokens

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(secret string, claims jwt.Claims) (string, error) {
	if secret == "" {
		return "", errors.New("JWT secret is required")
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ParseJWT(
	secret string,
	tokenString string,
	claims jwt.Claims,
	parserOptions ...jwt.ParserOption,
) error {
	if secret == "" {
		return errors.New("JWT secret is required")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(secret), nil
		},
		parserOptions...,
	)
	if err != nil {
		return err
	}
	if !token.Valid {
		return jwt.ErrTokenInvalidClaims
	}

	return nil
}
