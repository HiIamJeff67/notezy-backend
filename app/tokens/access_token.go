package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	exceptions "notezy-backend/app/exceptions"
	types "notezy-backend/shared/types"
)

/* ============================== Generate Token Functions ============================== */

func GenerateAccessToken(name string, email string, userAgent string) (*string, *exceptions.Exception) {
	claims := types.JWTClaims{
		Name:      name,
		Email:     email,
		UserAgent: userAgent,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_accessTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString(_jwtAccessTokenSecret)
	if err != nil {
		return nil, exceptions.Token.FailedToGenerateAccessToken().WithOrigin(err)
	}

	return &result, nil
}

/* ============================== Parse Token Functions ============================== */

func ParseAccessToken(tokenString string) (*types.JWTClaims, *exceptions.Exception) {
	accessToken, err := jwt.ParseWithClaims(
		tokenString,
		&types.JWTClaims{},
		func(token *jwt.Token) (any, error) { return _jwtAccessTokenSecret, nil },
	)
	if err != nil {
		return nil, exceptions.Token.FailedToParseAccessToken().WithOrigin(err)
	}

	if claims, ok := accessToken.Claims.(*types.JWTClaims); ok && accessToken.Valid {
		return claims, nil
	}

	return nil, exceptions.Token.FailedToParseAccessToken().WithOrigin(jwt.ErrTokenInvalidClaims)
}

/* ============================== Utility Functions ============================== */

func GetAccessTokenExpiresIn() time.Duration {
	return _accessTokenExpiresIn
}
