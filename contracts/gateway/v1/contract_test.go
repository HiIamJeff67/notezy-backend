package gatewaycontract

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestClientResponseDoesNotExposeAuthenticationTokens(t *testing.T) {
	response, err := json.Marshal(ClientResponse[struct{}]{
		Success: true,
		Data:    struct{}{},
		RefreshableTokens: &RefreshableTokens{
			NewCSRFToken: "csrf-token",
		},
	})
	if err != nil {
		t.Fatalf("marshal client response: %v", err)
	}

	body := string(response)
	if !strings.Contains(body, `"newCSRFToken":"csrf-token"`) {
		t.Fatalf("expected CSRF refresh metadata in response: %s", body)
	}
	if strings.Contains(body, "accessToken") || strings.Contains(body, "refreshToken") {
		t.Fatalf("authentication tokens must not be exposed in client response: %s", body)
	}
}

func TestClientResponseRoundTripsEmbeddedPublicId(t *testing.T) {
	publicId := uuid.New()
	original := ClientResponse[struct{}]{
		Success: true,
		Data:    struct{}{},
		Embedded: &EmbeddedResponse{
			PublicId: publicId,
		},
	}

	payload, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal client response: %v", err)
	}
	decoded := ClientResponse[struct{}]{}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal client response: %v", err)
	}
	if decoded.Embedded == nil || decoded.Embedded.PublicId != publicId {
		t.Fatalf("expected embedded public ID %s, got %#v", publicId, decoded.Embedded)
	}
}

func TestRequestRoundTripsTypedTokens(t *testing.T) {
	original := Request[struct{}]{
		Version:   Version,
		Operation: "auth.refresh",
		Tokens: Tokens{
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			CSRFToken:    "csrf-token",
		},
	}

	payload, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	decoded := Request[struct{}]{}
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal request: %v", err)
	}
	if decoded.Tokens != original.Tokens {
		t.Fatalf("expected typed tokens to round-trip, got %#v", decoded.Tokens)
	}
}

func TestClientRequestUsesDtoAsThePublicBody(t *testing.T) {
	type loginDto struct {
		Email string `json:"email"`
	}

	payload, err := json.Marshal(ClientRequest[loginDto]{
		Dto: loginDto{Email: "user@example.com"},
	})
	if err != nil {
		t.Fatalf("marshal client request: %v", err)
	}
	if string(payload) != `{"email":"user@example.com"}` {
		t.Fatalf("expected DTO-shaped public request, got %s", payload)
	}
}
