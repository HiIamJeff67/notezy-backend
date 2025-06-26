package cookies

import (
	"net/http"
	"time"

	shared "notezy-backend/shared"
	constants "notezy-backend/shared/constants"
)

var AccessToken = NewCookieHandler(
	shared.ValidCookieName_AccessToken,
	constants.BaseURL,
	time.Now().Add(constants.ExpirationTimeOfAccessToken),
	true,
	true,
	http.SameSiteLaxMode,
)
