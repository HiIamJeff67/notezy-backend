package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

/* ============================== Generate Tokens Functions ============================== */

func GenerateRefreshToken(id string, name string, email string, userAgent string) (*string, *exceptions.Exception) {
	claims := types.Claims{
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
		return nil, exceptions.Util.FailedToGenerateRefreshToken().WithError(err)
	}

	return &result, nil
}

/* ============================== Parse Tokens Functions ============================== */

func ParseRefreshToken(tokenString string) (*types.Claims, *exceptions.Exception) {
	refreshToken, err := jwt.ParseWithClaims(
		tokenString,
		&types.Claims{},
		func(token *jwt.Token) (any, error) { return _jwtRefreshTokenSecret, nil },
	)
	if err != nil {
		return nil, exceptions.Util.FailedToParseRefreshToken().WithError(err)
	}

	if claims, ok := refreshToken.Claims.(*types.Claims); ok && refreshToken.Valid {
		return claims, nil
	}

	return nil, exceptions.Util.FailedToParseRefreshToken().WithError(jwt.ErrTokenInvalidClaims)
}
