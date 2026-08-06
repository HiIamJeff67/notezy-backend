package middlewares

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtime-gateway/v1"
)

func TestDelegationMiddlewareAcceptsGatewayParticipantRequest(t *testing.T) {
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")
	t.Setenv("CORE_DELEGATION_AUDIENCE", "test-delegation-audience")
	t.Setenv("CORE_DELEGATION_ISSUER", "test-delegation-issuer")

	token, err := sharedtokens.GenerateDelegationToken(sharedtokens.DelegationTokenClaims{
		Actor:     "gateway",
		Operation: realtimegatewaycontract.GetBlockPackParticipantsOperation,
		RequestId: "request-id",
	})
	if err != nil {
		t.Fatalf("generate delegation token: %v", err)
	}

	payload, err := json.Marshal(realtimegatewaycontract.Request[realtimegatewaycontract.GetBlockPackParticipantsRequestDto]{
		Version:   realtimegatewaycontract.Version,
		Operation: realtimegatewaycontract.GetBlockPackParticipantsOperation,
		Metadata: realtimegatewaycontract.RequestMetadata{
			RequestId: "request-id",
		},
		Dto: realtimegatewaycontract.GetBlockPackParticipantsRequestDto{
			BlockPackId: uuid.New(),
		},
	})
	if err != nil {
		t.Fatalf("encode RealtimeGateway request: %v", err)
	}

	router := gin.New()
	router.POST(
		"/gateway/v1/block-pack-participants/get",
		DelegationMiddleware(realtimegatewaycontract.GetBlockPackParticipantsOperation),
		func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		},
	)

	request := httptest.NewRequest(http.MethodPost, "/gateway/v1/block-pack-participants/get", bytes.NewReader(payload))
	request.Header.Set("Authorization", "Bearer "+*token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-Id", "request-id")
	responseWriter := httptest.NewRecorder()
	router.ServeHTTP(responseWriter, request)
	if responseWriter.Code != http.StatusNoContent {
		t.Fatalf("expected delegated request to pass, got %d", responseWriter.Code)
	}
}
