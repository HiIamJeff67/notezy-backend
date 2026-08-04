package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthGoogleConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

func (c OAuthGoogleConfig) OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientId,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
