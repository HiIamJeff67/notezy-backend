package cookies

import (
	"net/http"
	"notezy-backend/global"
	"notezy-backend/global/constants"
	"time"
)

var RefreshToken = NewCookieHandler(
	global.ValidCookieName_RefreshToken,
	constants.BaseURL,
	time.Now().Add(constants.RefreshTokenExpirationTime),
	true,
	true,
	http.SameSiteStrictMode,
)
