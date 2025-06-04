package cookies

import (
	"net/http"
	"notezy-backend/global"
	"notezy-backend/global/constants"
	"time"
)

var AccessToken = NewCookieHandler(
	global.ValidCookieName_AccessToken,
	constants.BaseURL,
	time.Now().Add(constants.AccessTokenExpirationTime),
	true,
	true,
	http.SameSiteLaxMode,
)
