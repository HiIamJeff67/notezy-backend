package cookies

import (
	"net/http"
	"time"

	shared "notezy-backend/shared"
	constants "notezy-backend/shared/constants"
)

var RefreshToken = NewCookieHandler(
	shared.ValidCookieName_RefreshToken,                    // name
	constants.BaseURL,                                      // path
	time.Now().Add(constants.ExpirationTimeOfRefreshToken), // expires
	true,                    // secure
	true,                    // httpOnly
	http.SameSiteStrictMode, // sameSite
)
