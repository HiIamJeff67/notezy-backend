package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtimegateway/v1"
)

func TestRealtimeGatewayClientRequestsVersionedParticipantSnapshot(t *testing.T) {
	t.Setenv("CORE_DELEGATION_SECRET", "test-delegation-secret")
	t.Setenv("CORE_DELEGATION_AUDIENCE", "test-delegation-audience")
	t.Setenv("CORE_DELEGATION_ISSUER", "test-delegation-issuer")

	blockPackId := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			t.Fatalf("expected POST request, got %s", request.Method)
		}
		if request.URL.Path != getBlockPackParticipantsPath {
			t.Fatalf("expected participant path, got %s", request.URL.Path)
		}
		if request.Header.Get("Authorization") == "" {
			t.Fatal("expected delegation credential")
		}

		requestEnvelope := realtimegatewaycontract.Request[realtimegatewaycontract.GetBlockPackParticipantsRequestDto]{}
		if err := json.NewDecoder(request.Body).Decode(&requestEnvelope); err != nil {
			t.Fatalf("decode RealtimeGateway request: %v", err)
		}
		if requestEnvelope.Version != realtimegatewaycontract.Version ||
			requestEnvelope.Operation != realtimegatewaycontract.GetBlockPackParticipantsOperation ||
			requestEnvelope.Dto.BlockPackId != blockPackId {
			t.Fatal("expected versioned RealtimeGateway participant request")
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(responseWriter).Encode(realtimegatewaycontract.Response[realtimegatewaycontract.GetBlockPackParticipantsResponseDto]{
			Version: realtimegatewaycontract.Version,
			Metadata: realtimegatewaycontract.ResponseMetadata{
				RequestId:   requestEnvelope.Metadata.RequestId,
				RespondedAt: time.Now(),
			},
			Data: realtimegatewaycontract.GetBlockPackParticipantsResponseDto{
				Participants: []realtimegatewaycontract.BlockPackParticipantResponseDto{
					{
						UserPublicId:      uuid.New(),
						ChannelPermission: "write",
						ConnectionCount:   1,
					},
				},
			},
		}); err != nil {
			t.Fatalf("encode RealtimeGateway response: %v", err)
		}
	}))
	defer server.Close()

	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("X-Request-Id", "request-id")
	ctx.Request.Header.Set("Traceparent", "00-trace")

	client := NewRealtimeGatewayClient(server.URL, time.Second)
	responseDto, exception := client.GetBlockPackParticipants(
		ctx,
		&realtimegatewaycontract.GetBlockPackParticipantsRequestDto{
			BlockPackId: blockPackId,
		},
	)
	if exception != nil {
		t.Fatalf("get RealtimeGateway participants: %v", exception)
	}
	if len(responseDto.Participants) != 1 || responseDto.Participants[0].ChannelPermission != "write" {
		t.Fatal("expected RealtimeGateway participant snapshot")
	}
}
