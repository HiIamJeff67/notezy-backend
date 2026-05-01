package cookies

import (
	"net/http"

	constants "notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

var AccessTokenCookieHandler = NewCookieHandler(
	types.ValidCookieName_AccessToken,     // name
	"/",                                   // path
	constants.ExpirationTimeOfAccessToken, // duration
	constants.Mode == types.ModeType_Production, // secure (set to true only if is on the production)
	true,                 // httpOnly
	http.SameSiteLaxMode, // sameSite
)

// Note: make sure the path should start with "/" because we want this work at the all the subpath from "/"
