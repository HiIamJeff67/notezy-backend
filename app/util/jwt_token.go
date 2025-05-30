package util

import (
	"time"

	exceptions "notezy-backend/app/exceptions"
	global "notezy-backend/global"
	types "notezy-backend/global/types"

	"github.com/golang-jwt/jwt/v5"
)

/* ============================== Get The Secret Key Storing in the Environment ============================== */
var _jwtAccessTokenSecret []byte
var _jwtRefreshTokenSecret []byte

const (
	_accessTokenExpiresIn  time.Duration = 30 * time.Minute
	_refreshTokenExpiresIn time.Duration = 7 * 24 * time.Hour
)

func init() {
	accessTokenSecretKey := global.GetEnv("JWT_ACCESS_TOKEN_SECRET_KEY", "")
	refreshTokenSecretKey := global.GetEnv("JWT_REFRESH_TOKEN_SECRET_KEY", "")
	if accessTokenSecretKey == "" {
		exceptions.Util.AccessTokenSecretKeyNotFound()
	}
	if refreshTokenSecretKey == "" {
		exceptions.Util.RefreshTokenSecretKeyNotFound()
	}
	_jwtAccessTokenSecret = []byte(accessTokenSecretKey)
	_jwtRefreshTokenSecret = []byte(refreshTokenSecretKey)
}

/* ============================== Generate Tokens Functions ============================== */
func GenerateAccessToken(id string, name string, email string) (string, *exceptions.Exception) {
	claims := types.Claims{
		Id:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_accessTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString(_jwtAccessTokenSecret)
	if err != nil {
		return "", exceptions.Util.FailedToGenerateAccessToken().WithError(err)
	}

	return result, nil
}

func GenerateRefreshToken(id string, name string, email string) (string, *exceptions.Exception) {
	claims := types.Claims{
		Id:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_refreshTokenExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString(_jwtRefreshTokenSecret)
	if err != nil {
		return "", exceptions.Util.FailedToGenerateRefreshToken().WithError(err)
	}

	return result, nil
}

/* ============================== Parse Tokens Functions ============================== */
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
