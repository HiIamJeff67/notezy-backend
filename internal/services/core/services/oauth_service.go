package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type OAuthServiceInterface interface {
	GetGoogleUserInfo(ctx context.Context, authenticationCode string) (*googleUserInfo, *exceptions.Exception)
}

type OAuthService struct {
	oauthGoogleConfig *oauth2.Config
}

func NewOAuthService(oauthGoogleConfig *oauth2.Config) OAuthServiceInterface {
	return &OAuthService{
		oauthGoogleConfig: oauthGoogleConfig,
	}
}

/* ============================== Service Methods for OAuth ============================== */

type googleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verifiedEmail"`
	Name          string `json:"name"`
	GivenName     string `json:"givenName"`
	FamilyName    string `json:"familyName"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func (s *OAuthService) GetGoogleUserInfo(
	ctx context.Context, authenticationCode string,
) (*googleUserInfo, *exceptions.Exception) {
	token, err := s.oauthGoogleConfig.Exchange(ctx, authenticationCode)
	if err != nil {
		return nil, exceptions.New(
			"TokenExchangeFailed",
			"OAuth",
			"GetGoogleUserInfo",
			"Failed to exchange the OAuth token",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	client := s.oauthGoogleConfig.Client(ctx, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, exceptions.New(
			"InvalidAuthenticationCode",
			"OAuth",
			"GetGoogleUserInfo",
			"Authentication code is invalid",
			http.StatusBadRequest,
		).WithOrigin(err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, exceptions.New(
			"ResponseReadFailed",
			"OAuth",
			"GetGoogleUserInfo",
			"Failed to read the OAuth provider response",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	var userInfo googleUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, exceptions.New(
			"InvalidResponse",
			"OAuth",
			"GetGoogleUserInfo",
			"OAuth provider response is invalid",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	return &userInfo, nil
}
