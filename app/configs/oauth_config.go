package configs

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	util "notezy-backend/app/util"
)

var OAuthGoogleConfig = &oauth2.Config{
	ClientID:     util.GetEnv("OAUTH_GOOGLE_CLIENT_ID", ""),
	ClientSecret: util.GetEnv("OAUTH_GOOGLE_CLIENT_SECRET", ""),
	RedirectURL:  util.GetEnv("OAUTH_GOOGLE_REDIRECT_URL", ""),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}
