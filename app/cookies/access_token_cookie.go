package cookies

import (
	"net/http"

	shared "notezy-backend/shared"
	constants "notezy-backend/shared/constants"
)

var AccessToken = NewCookieHandler(
	shared.ValidCookieName_AccessToken,    // name
	"/"+constants.CurrentBaseURL,          // path
	constants.ExpirationTimeOfAccessToken, // duration
	true,                                  // secure
	true,                                  // httpOnly
	http.SameSiteLaxMode,                  // sameSite
)

// Note: make sure the path should start with "/" because we want this work at the all the subpath from constants.CurrentBaseURL
