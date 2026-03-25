package services

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"

	dtos "notezy-backend/app/dtos"
	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	traces "notezy-backend/app/traces"
)

type OAuthServiceInterface interface {
	GetGoogleUserInfo(ctx context.Context, authenticationCode string) (*dtos.GoogleUserInfoDto, *exceptions.Exception)
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

func (s *OAuthService) GetGoogleUserInfo(
	ctx context.Context, authenticationCode string,
) (*dtos.GoogleUserInfoDto, *exceptions.Exception) {
	logs.Info(traces.GetTrace(0).FileLineString(), "Client Id: ", s.oauthGoogleConfig.ClientID)
	logs.Info(traces.GetTrace(0).FileLineString(), "Client Secret: ", s.oauthGoogleConfig.ClientSecret)
	logs.Info(traces.GetTrace(0).FileLineString(), "Redirect URL: ", s.oauthGoogleConfig.RedirectURL)
	token, err := s.oauthGoogleConfig.Exchange(ctx, authenticationCode)
	if err != nil {
		return nil, exceptions.OAuth.FailedToExchangeToken(authenticationCode).WithError(err)
	}

	client := s.oauthGoogleConfig.Client(ctx, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, exceptions.OAuth.InvalidAuthenticationCode(authenticationCode).WithError(err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, exceptions.OAuth.FailedToParseResposneFromOAuthThirdParty("google").WithError(err)
	}

	var userInfo dtos.GoogleUserInfoDto
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, exceptions.OAuth.InvalidDto().WithError(err)
	}

	return &userInfo, nil
}
