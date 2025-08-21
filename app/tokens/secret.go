package tokens

import (
	"time"

	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
)

var _jwtAccessTokenSecret []byte
var _jwtRefreshTokenSecret []byte

const (
	_accessTokenExpiresIn  time.Duration = 30 * time.Minute
	_refreshTokenExpiresIn time.Duration = 7 * 24 * time.Hour
)

func init() {
	accessTokenSecretKey := util.GetEnv("JWT_ACCESS_TOKEN_SECRET_KEY", "")
	refreshTokenSecretKey := util.GetEnv("JWT_REFRESH_TOKEN_SECRET_KEY", "")
	if accessTokenSecretKey == "" {
		exceptions.Util.AccessTokenSecretKeyNotFound()
	}
	if refreshTokenSecretKey == "" {
		exceptions.Util.RefreshTokenSecretKeyNotFound()
	}
	_jwtAccessTokenSecret = []byte(accessTokenSecretKey)
	_jwtRefreshTokenSecret = []byte(refreshTokenSecretKey)
}
