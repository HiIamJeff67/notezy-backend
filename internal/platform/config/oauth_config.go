package configs

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var OAuthGoogleConfig = &oauth2.Config{
	ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("OAUTH_GOOGLE_REDIRECT_URL"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

var OAuthPaypalConfig = &oauth2.Config{}
