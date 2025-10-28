package cookies

import (
	"net/http"

	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

var RefreshToken = NewCookieHandler(
	types.ValidCookieName_RefreshToken,     // name
	"/",                                    // path
	constants.ExpirationTimeOfRefreshToken, // duration
	true,                                   // secure
	true,                                   // httpOnly
	http.SameSiteStrictMode,                // sameSite
)

// Note: make sure the path should start with "/" because we want this work at the all the subpath from "/"
