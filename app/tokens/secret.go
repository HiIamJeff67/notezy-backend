package tokens

import (
	"notezy-backend/app/exceptions"
	"notezy-backend/app/shared"
	"time"
)

var _jwtAccessTokenSecret []byte
var _jwtRefreshTokenSecret []byte

const (
	_accessTokenExpiresIn  time.Duration = 30 * time.Minute
	_refreshTokenExpiresIn time.Duration = 7 * 24 * time.Hour
)

func init() {
	accessTokenSecretKey := shared.GetEnv("JWT_ACCESS_TOKEN_SECRET_KEY", "")
	refreshTokenSecretKey := shared.GetEnv("JWT_REFRESH_TOKEN_SECRET_KEY", "")
	if accessTokenSecretKey == "" {
		exceptions.Util.AccessTokenSecretKeyNotFound()
	}
	if refreshTokenSecretKey == "" {
		exceptions.Util.RefreshTokenSecretKeyNotFound()
	}
	_jwtAccessTokenSecret = []byte(accessTokenSecretKey)
	_jwtRefreshTokenSecret = []byte(refreshTokenSecretKey)
}
