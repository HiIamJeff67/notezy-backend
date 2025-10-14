package tokens

import (
	"time"

	exceptions "notezy-backend/app/exceptions"
	util "notezy-backend/app/util"
)

var _jwtAccessTokenSecret []byte
var _jwtRefreshTokenSecret []byte
var _csrfTokenSecret []byte

const (
	_accessTokenExpiresIn  time.Duration = 30 * time.Minute
	_refreshTokenExpiresIn time.Duration = 7 * 24 * time.Hour
	_csrfTokenExpiresIn    time.Duration = 7 * 24 * time.Hour
)

const (
	_csrfTokenLength = 32
)

func init() {
	accessTokenSecretKey := util.GetEnv("JWT_ACCESS_TOKEN_SECRET_KEY", "")
	refreshTokenSecretKey := util.GetEnv("JWT_REFRESH_TOKEN_SECRET_KEY", "")
	csrfTokenSecretKey := util.GetEnv("CSRF_TOKEN_SECRET_KEY", "")

	if accessTokenSecretKey == "" {
		exceptions.Token.AccessTokenSecretKeyNotFound()
	}
	if refreshTokenSecretKey == "" {
		exceptions.Token.RefreshTokenSecretKeyNotFound()
	}
	if csrfTokenSecretKey == "" {
		exceptions.Token.CSRFTokenSecretKeyNotFound()
	}

	_jwtAccessTokenSecret = []byte(accessTokenSecretKey)
	_jwtRefreshTokenSecret = []byte(refreshTokenSecretKey)
	_csrfTokenSecret = []byte(csrfTokenSecretKey)
}
