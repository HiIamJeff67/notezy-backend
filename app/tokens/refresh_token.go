package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

/* ============================== Generate Token Functions ============================== */

func GenerateRefreshToken(id string, name string, email string, userAgent string) (*string, *exceptions.Exception) {
	claims := types.JWTClaims{
		Id:        id,
		Name:      name,
		Email:     email,
		UserAgent: userAgent,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_refreshTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString(_jwtRefreshTokenSecret)
	if err != nil {
		return nil, exceptions.Token.FailedToGenerateRefreshToken().WithError(err)
	}

	return &result, nil
}

/* ============================== Parse Token Functions ============================== */

func ParseRefreshToken(tokenString string) (*types.JWTClaims, *exceptions.Exception) {
	refreshToken, err := jwt.ParseWithClaims(
		tokenString,
		&types.JWTClaims{},
		func(token *jwt.Token) (any, error) { return _jwtRefreshTokenSecret, nil },
	)
	if err != nil {
		return nil, exceptions.Token.FailedToParseRefreshToken().WithError(err)
	}

	if claims, ok := refreshToken.Claims.(*types.JWTClaims); ok && refreshToken.Valid {
		return claims, nil
	}

	return nil, exceptions.Token.FailedToParseRefreshToken().WithError(jwt.ErrTokenInvalidClaims)
}

/* ============================== Utility Functions ============================== */

func GetRefreshTokenExpiresIn() time.Duration {
	return _refreshTokenExpiresIn
}
