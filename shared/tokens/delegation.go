package tokens

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type DelegationTokenClaims struct {
	Actor              string   `json:"actor"`
	GatewaySource      string   `json:"gatewaySource,omitempty"`
	AuthMethod         string   `json:"authMethod,omitempty"`
	ApiKeyId           string   `json:"apiKeyId,omitempty"`
	UserSubject        string   `json:"userSubject,omitempty"`
	AllowedPermissions []string `json:"allowedPermissions"`
	Operation          string   `json:"operation"`
	RequestId          string   `json:"requestId"`
	jwt.RegisteredClaims
}

const (
	GatewaySourceClient = "client"
	GatewaySourceAPI    = "api"

	AuthMethodJWT    = "jwt"
	AuthMethodAPIKey = "api-key"
)

func GenerateDelegationToken(claims DelegationTokenClaims) (*string, error) {
	if claims.Actor == "" || claims.Operation == "" || claims.RequestId == "" {
		return nil, errors.New("delegation token claims are invalid")
	}
	if claims.GatewaySource != "" && claims.GatewaySource != GatewaySourceClient && claims.GatewaySource != GatewaySourceAPI {
		return nil, errors.New("delegation gateway source is invalid")
	}
	if claims.AuthMethod != "" && claims.AuthMethod != AuthMethodJWT && claims.AuthMethod != AuthMethodAPIKey {
		return nil, errors.New("delegation auth method is invalid")
	}
	if claims.GatewaySource == GatewaySourceAPI && claims.AuthMethod != AuthMethodAPIKey {
		return nil, errors.New("api gateway delegation requires api key authentication")
	}
	if claims.GatewaySource == GatewaySourceClient && claims.AuthMethod == AuthMethodAPIKey {
		return nil, errors.New("client gateway delegation cannot use api key authentication")
	}
	if claims.UserSubject != "" {
		if _, err := uuid.Parse(claims.UserSubject); err != nil {
			return nil, errors.New("delegation user subject is invalid")
		}
	}
	secret := os.Getenv("CORE_DELEGATION_SECRET")
	if secret == "" {
		return nil, errors.New("delegation token secret is required")
	}
	now := time.Now()
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{os.Getenv("CORE_DELEGATION_AUDIENCE")},
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
		ID:        uuid.NewString(),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    os.Getenv("CORE_DELEGATION_ISSUER"),
		Subject:   claims.UserSubject,
	}

	token, err := SignJWT(secret, claims)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func ParseDelegationToken(tokenString string) (*DelegationTokenClaims, error) {
	claims := &DelegationTokenClaims{}
	err := ParseJWT(
		os.Getenv("CORE_DELEGATION_SECRET"),
		tokenString,
		claims,
		jwt.WithAudience(os.Getenv("CORE_DELEGATION_AUDIENCE")),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(os.Getenv("CORE_DELEGATION_ISSUER")),
	)
	if err != nil {
		return nil, err
	}
	if claims.Actor == "" || claims.Operation == "" || claims.RequestId == "" {
		return nil, errors.New("delegation token claims are invalid")
	}
	if claims.GatewaySource != "" && claims.GatewaySource != GatewaySourceClient && claims.GatewaySource != GatewaySourceAPI {
		return nil, errors.New("delegation gateway source is invalid")
	}
	if claims.AuthMethod != "" && claims.AuthMethod != AuthMethodJWT && claims.AuthMethod != AuthMethodAPIKey {
		return nil, errors.New("delegation auth method is invalid")
	}
	if claims.GatewaySource == GatewaySourceAPI && claims.AuthMethod != AuthMethodAPIKey {
		return nil, errors.New("api gateway delegation requires api key authentication")
	}
	if claims.GatewaySource == GatewaySourceClient && claims.AuthMethod == AuthMethodAPIKey {
		return nil, errors.New("client gateway delegation cannot use api key authentication")
	}
	if claims.Subject != claims.UserSubject {
		return nil, errors.New("delegation token claims are invalid")
	}

	return claims, nil
}
