package tokens

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/shared/types"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/* ============================== Generate Tokens Function ============================== */

func GenerateAccessToken(id string, name string, email string, userAgent string) (*string, *exceptions.Exception) {
	claims := types.Claims{
		Id:        id,
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
		return nil, exceptions.Util.FailedToGenerateAccessToken().WithError(err)
	}

	return &result, nil
}

/* ============================== Parse Tokens Function ============================== */

func ParseAccessToken(tokenString string) (*types.Claims, *exceptions.Exception) {
	accessToken, err := jwt.ParseWithClaims(
		tokenString,
		&types.Claims{},
		func(token *jwt.Token) (any, error) { return _jwtAccessTokenSecret, nil },
	)
	if err != nil {
		return nil, exceptions.Util.FailedToParseAccessToken().WithError(err)
	}

	if claims, ok := accessToken.Claims.(*types.Claims); ok && accessToken.Valid {
		return claims, nil
	}

	return nil, exceptions.Util.FailedToParseAccessToken().WithError(jwt.ErrTokenInvalidClaims)
}
