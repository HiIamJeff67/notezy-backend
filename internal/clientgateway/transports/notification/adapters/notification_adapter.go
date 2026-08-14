package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

	sharedcontexts "github.com/HiIamJeff67/notezy-backend/shared/lib/contexts"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/shared/tokens"

	gatewaycontexts "github.com/HiIamJeff67/notezy-backend/internal/clientgateway/contexts"
)

type NotificationAdapter struct {
	baseURL    string
	httpClient *http.Client
}

func NewNotificationAdapter(baseURL string, timeout time.Duration) *NotificationAdapter {
	return &NotificationAdapter{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func CallSecurly[RequestDto any, ResponseDto any](
	ctx *gin.Context,
	client *NotificationAdapter,
	requestDto *RequestDto,
	operation string,
	path string,
) (*gatewaycontract.Response[ResponseDto], *exceptions.Exception) {
	if client == nil {
		return nil, exceptions.New(
			"NotificationAdapterRequired",
			"Gateway",
			operation,
			"The Notification service adapter is required",
			http.StatusInternalServerError,
			true,
		)
	}
	if requestDto == nil {
		return nil, exceptions.New(
			"InvalidRequest",
			"Gateway",
			operation,
			"The Notification service request DTO is required",
			http.StatusBadRequest,
		)
	}

	userSubject, exception := gatewaycontexts.GetAndConvertContextFieldToUUID(
		ctx,
		sharedcontexts.ContextFieldName_User_PublicId,
	)
	if exception != nil {
		return nil, exception
	}
	if userSubject == nil || *userSubject == uuid.Nil {
		return nil, exceptions.New(
			"ContextFieldInvalid",
			"Gateway",
			operation,
			"A valid user subject is required for a secure Notification service call",
			http.StatusInternalServerError,
			true,
		)
	}

	requestId := ctx.GetHeader("X-Request-Id")
	if requestId == "" {
		requestId = uuid.NewString()
	}
	delegationToken, err := sharedtokens.GenerateDelegationToken(sharedtokens.DelegationTokenClaims{
		Actor:       "gateway",
		UserSubject: userSubject.String(),
		Operation:   operation,
		RequestId:   requestId,
	})
	if err != nil {
		return nil, exceptions.New(
			"NotificationDelegationFailed",
			"Gateway",
			operation,
			"Failed to issue a Notification service delegation token",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	request := &gatewaycontract.Request[RequestDto]{
		Operation: operation,
		Metadata: gatewaycontract.RequestMetadata{
			RequestId:      requestId,
			TraceParent:    ctx.GetHeader("Traceparent"),
			IdempotencyKey: ctx.GetHeader("Idempotency-Key"),
		},
		Dto: *requestDto,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return nil, exceptions.New(
			"NotificationRequestEncodingFailed",
			"Gateway",
			operation,
			"Failed to encode the Notification service request",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx.Request.Context(),
		http.MethodPost,
		client.baseURL+"/"+strings.TrimLeft(path, "/"),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, exceptions.New(
			"NotificationRequestCreationFailed",
			"Gateway",
			operation,
			"Failed to create the Notification service request",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+*delegationToken)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("X-Request-Id", requestId)
	if request.Metadata.TraceParent != "" {
		httpRequest.Header.Set("Traceparent", request.Metadata.TraceParent)
	}
	if request.Metadata.IdempotencyKey != "" {
		httpRequest.Header.Set("Idempotency-Key", request.Metadata.IdempotencyKey)
	}
	for _, header := range []string{"User-Agent", "X-Real-IP", "X-Forwarded-For"} {
		if value := ctx.GetHeader(header); value != "" {
			httpRequest.Header.Set(header, value)
		}
	}

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, exceptions.New(
			"NotificationRequestFailed",
			"Gateway",
			operation,
			"Failed to communicate with the Notification service",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	defer httpResponse.Body.Close()
	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, exceptions.New(
			"NotificationResponseReadFailed",
			"Gateway",
			operation,
			"Failed to read the Notification service response",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		response := &gatewaycontract.Response[ResponseDto]{}
		if err := json.Unmarshal(responseBody, response); err == nil && response.Exception != nil {
			return nil, response.Exception.Clone(httpResponse.StatusCode)
		}
		return nil, exceptions.New(
			"NotificationResponseFailed",
			"Gateway",
			operation,
			"The Notification service returned an unsuccessful response",
			httpResponse.StatusCode,
			true,
		).WithOrigin(fmt.Errorf("status %d: %s", httpResponse.StatusCode, strings.TrimSpace(string(responseBody))))
	}
	response := &gatewaycontract.Response[ResponseDto]{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, exceptions.New(
			"NotificationResponseDecodingFailed",
			"Gateway",
			operation,
			"Failed to decode the Notification service response",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	if response.Version != gatewaycontract.Version {
		return nil, exceptions.New(
			"NotificationResponseVersionInvalid",
			"Gateway",
			operation,
			"The Notification service response uses an unsupported version",
			http.StatusInternalServerError,
			true,
		)
	}
	if response.Metadata.RequestId != requestId {
		return nil, exceptions.New(
			"NotificationResponseRequestIdInvalid",
			"Gateway",
			operation,
			"The Notification service response does not match the request",
			http.StatusInternalServerError,
			true,
		)
	}
	return response, nil
}
