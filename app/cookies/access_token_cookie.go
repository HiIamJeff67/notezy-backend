package cookies

import (
	"net/http"
	"time"

	shared "notezy-backend/shared"
	constants "notezy-backend/shared/constants"
)

var AccessToken = NewCookieHandler(
	shared.ValidCookieName_AccessToken,                    // name
	constants.BaseURL,                                     // path
	time.Now().Add(constants.ExpirationTimeOfAccessToken), // expires
	true,                 // secure
	true,                 // httpOnly
	http.SameSiteLaxMode, // sameSite
)
