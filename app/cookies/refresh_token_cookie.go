package cookies

import (
	"net/http"

	shared "notezy-backend/shared"
	constants "notezy-backend/shared/constants"
)

var RefreshToken = NewCookieHandler(
	shared.ValidCookieName_RefreshToken,    // name
	"/"+constants.CurrentBaseURL,           // path
	constants.ExpirationTimeOfRefreshToken, // duration
	true,                                   // secure
	true,                                   // httpOnly
	http.SameSiteStrictMode,                // sameSite
)

// Note: make sure the path should start with "/" because we want this work at the all the subpath from constants.CurrentBaseURL
