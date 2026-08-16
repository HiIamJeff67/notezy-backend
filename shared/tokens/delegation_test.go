package tokens

import "testing"

func TestDelegationTokenRoundTrip(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notegic-core-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notegic-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")

	token, err := GenerateDelegationToken(DelegationTokenClaims{
		Actor:              "gateway",
		UserSubject:        "83bdeac1-02de-42fe-a7a8-4e1a83174866",
		AllowedPermissions: []string{"Read", "Write"},
		Operation:          "root-shelves.update",
		RequestId:          "request-id",
	})
	if err != nil {
		t.Fatalf("generate delegation token: %v", err)
	}

	claims, err := ParseDelegationToken(*token)
	if err != nil {
		t.Fatalf("parse delegation token: %v", err)
	}
	if claims.Actor != "gateway" || claims.UserSubject != "83bdeac1-02de-42fe-a7a8-4e1a83174866" ||
		claims.Operation != "root-shelves.update" ||
		claims.RequestId != "request-id" ||
		len(claims.AllowedPermissions) != 2 {
		t.Fatalf("unexpected delegation token claims: %#v", claims)
	}
}

func TestDelegationTokenWithoutUserSubjectRoundTrip(t *testing.T) {
	t.Setenv("CORE_DELEGATION_AUDIENCE", "notegic-core-test")
	t.Setenv("CORE_DELEGATION_ISSUER", "notegic-gateway-test")
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")

	token, err := GenerateDelegationToken(DelegationTokenClaims{
		Actor:     "gateway",
		Operation: "auth.login",
		RequestId: "request-id",
	})
	if err != nil {
		t.Fatalf("generate delegation token without user subject: %v", err)
	}

	claims, err := ParseDelegationToken(*token)
	if err != nil {
		t.Fatalf("parse delegation token without user subject: %v", err)
	}
	if claims.Actor != "gateway" || claims.UserSubject != "" || claims.Subject != "" {
		t.Fatalf("unexpected delegation claims without user subject: %#v", claims)
	}
}

func TestDelegationTokenCarriesGatewaySourceAndAPIKeyMetadata(t *testing.T) {
	t.Setenv("CORE_DELEGATION_SECRET", "delegation-secret")
	t.Setenv("CORE_DELEGATION_AUDIENCE", "core")
	t.Setenv("CORE_DELEGATION_ISSUER", "gateway")
	token, err := GenerateDelegationToken(DelegationTokenClaims{
		Actor:         "gateway",
		GatewaySource: GatewaySourceAPI,
		AuthMethod:    AuthMethodAPIKey,
		ApiKeyId:      "api-key-id",
		Operation:     "get-user",
		RequestId:     "request-id",
	})
	if err != nil {
		t.Fatalf("generate API delegation token: %v", err)
	}
	claims, err := ParseDelegationToken(*token)
	if err != nil {
		t.Fatalf("parse API delegation token: %v", err)
	}
	if claims.GatewaySource != GatewaySourceAPI || claims.AuthMethod != AuthMethodAPIKey || claims.ApiKeyId != "api-key-id" {
		t.Fatalf("unexpected API delegation metadata: %+v", claims)
	}
}
