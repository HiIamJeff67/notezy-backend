package adapters

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	realtimegatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/realtime-gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"
)

const getBlockPackParticipantsPath = "/gateway/v1/block-pack-participants/get"

type RealtimeGatewayClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewRealtimeGatewayClient(baseURL string, timeout time.Duration) *RealtimeGatewayClient {
	return &RealtimeGatewayClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *RealtimeGatewayClient) GetBlockPackParticipants(
	ctx *gin.Context,
	requestDto *realtimegatewaycontract.GetBlockPackParticipantsRequestDto,
) (*realtimegatewaycontract.GetBlockPackParticipantsResponseDto, *exceptions.Exception) {
	if c == nil || c.httpClient == nil || c.baseURL == "" {
		return nil, exceptions.New(
			"RealtimeGatewayClientRequired",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"The RealtimeGateway service client is required",
			http.StatusInternalServerError,
			true,
		)
	}
	if requestDto == nil || requestDto.BlockPackId == uuid.Nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"The RealtimeGateway participant request is invalid",
			http.StatusBadRequest,
		)
	}

	requestId := ctx.GetHeader("X-Request-Id")
	if requestId == "" {
		requestId = uuid.NewString()
	}

	delegationToken, err := sharedtokens.GenerateDelegationToken(sharedtokens.DelegationTokenClaims{
		Actor:     "gateway",
		Operation: realtimegatewaycontract.GetBlockPackParticipantsOperation,
		RequestId: requestId,
	})
	if err != nil {
		return nil, exceptions.New(
			"RealtimeGatewayDelegationFailed",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"Failed to communicate with the RealtimeGateway service",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	payload, err := json.Marshal(realtimegatewaycontract.Request[realtimegatewaycontract.GetBlockPackParticipantsRequestDto]{
		Version:   realtimegatewaycontract.Version,
		Operation: realtimegatewaycontract.GetBlockPackParticipantsOperation,
		Metadata: realtimegatewaycontract.RequestMetadata{
			RequestId:   requestId,
			TraceParent: ctx.GetHeader("Traceparent"),
		},
		Dto: *requestDto,
	})
	if err != nil {
		return nil, exceptions.New(
			"RealtimeGatewayRequestEncodingFailed",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"Failed to encode the RealtimeGateway participant request",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	request, err := http.NewRequestWithContext(
		ctx.Request.Context(),
		http.MethodPost,
		c.baseURL+getBlockPackParticipantsPath,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, exceptions.New(
			"RealtimeGatewayRequestCreationFailed",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"Failed to create the RealtimeGateway participant request",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	request.Header.Set("Authorization", "Bearer "+*delegationToken)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Request-Id", requestId)
	if traceParent := ctx.GetHeader("Traceparent"); traceParent != "" {
		request.Header.Set("Traceparent", traceParent)
	}

	httpResponse, err := c.httpClient.Do(request)
	if err != nil {
		return nil, exceptions.New(
			"RealtimeGatewayRequestFailed",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"Failed to communicate with the RealtimeGateway service",
			http.StatusServiceUnavailable,
			true,
		).WithOrigin(err)
	}
	defer httpResponse.Body.Close()

	responsePayload, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, exceptions.New(
			"RealtimeGatewayResponseReadFailed",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"Failed to read the RealtimeGateway participant response",
			http.StatusServiceUnavailable,
			true,
		).WithOrigin(err)
	}

	response := &realtimegatewaycontract.Response[realtimegatewaycontract.GetBlockPackParticipantsResponseDto]{}
	if err := json.Unmarshal(responsePayload, response); err != nil {
		return nil, exceptions.New(
			"RealtimeGatewayResponseDecodingFailed",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"Failed to decode the RealtimeGateway participant response",
			http.StatusServiceUnavailable,
			true,
		).WithOrigin(err)
	}
	if response.Version != realtimegatewaycontract.Version || response.Metadata.RequestId != requestId {
		return nil, exceptions.New(
			"RealtimeGatewayResponseInvalid",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"The RealtimeGateway participant response is invalid",
			http.StatusServiceUnavailable,
			true,
		)
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		if response.Exception != nil {
			return nil, response.Exception.Clone(httpResponse.StatusCode)
		}

		return nil, exceptions.New(
			"RealtimeGatewayResponseFailed",
			"Gateway",
			realtimegatewaycontract.GetBlockPackParticipantsOperation,
			"The RealtimeGateway service returned an unsuccessful response",
			http.StatusServiceUnavailable,
			true,
		)
	}

	return &response.Data, nil
}
