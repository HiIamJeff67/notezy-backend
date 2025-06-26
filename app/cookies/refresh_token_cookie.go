package cookies

import (
	"net/http"
	"time"

	shared "notezy-backend/shared"
	constants "notezy-backend/shared/constants"
)

var RefreshToken = NewCookieHandler(
	shared.ValidCookieName_RefreshToken,
	constants.BaseURL,
	time.Now().Add(constants.ExpirationTimeOfRefreshToken),
	true,
	true,
	http.SameSiteStrictMode,
)
