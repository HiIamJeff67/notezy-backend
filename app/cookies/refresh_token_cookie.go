package cookies

import (
	"net/http"

	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

var RefreshTokenCookieHandler = NewCookieHandler(
	types.ValidCookieName_RefreshToken,     // name
	"/",                                    // path
	constants.ExpirationTimeOfRefreshToken, // duration
	constants.Mode == types.ModeType_Production, // secure (set to true only if is on the production)
	true,                    // httpOnly
	http.SameSiteStrictMode, // sameSite
)

// Note: make sure the path should start with "/" because we want this work at the all the subpath from "/"
