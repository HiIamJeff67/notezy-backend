package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	emaildto "github.com/HiIamJeff67/notezy-backend/contracts/email/v1"
	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
)

type ClientInterface interface {
	SendWelcomeEmail(ctx context.Context, requestDto emaildto.SendWelcomeEmailRequestDto) *exceptions.Exception
	SendValidationEmail(ctx context.Context, requestDto emaildto.SendValidationEmailRequestDto) *exceptions.Exception
	SendSecurityAlertEmail(ctx context.Context, requestDto emaildto.SendSecurityAlertEmailRequestDto) *exceptions.Exception
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) ClientInterface {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) SendWelcomeEmail(
	ctx context.Context,
	requestDto emaildto.SendWelcomeEmailRequestDto,
) *exceptions.Exception {
	return c.send(
		ctx,
		emaildto.SendWelcomeEmailOperation,
		"/email/v1/send/welcome",
		requestDto,
	)
}

func (c *Client) SendValidationEmail(
	ctx context.Context,
	requestDto emaildto.SendValidationEmailRequestDto,
) *exceptions.Exception {
	return c.send(
		ctx,
		emaildto.SendValidationEmailOperation,
		"/email/v1/send/validation",
		requestDto,
	)
}

func (c *Client) SendSecurityAlertEmail(
	ctx context.Context,
	requestDto emaildto.SendSecurityAlertEmailRequestDto,
) *exceptions.Exception {
	return c.send(
		ctx,
		emaildto.SendSecurityAlertEmailOperation,
		"/email/v1/send/security-alert",
		requestDto,
	)
}

func (c *Client) send(
	ctx context.Context,
	operation string,
	path string,
	requestDto any,
) *exceptions.Exception {
	request := gatewaycontract.Request[any]{
		Version:   gatewaycontract.Version,
		Operation: operation,
		Metadata: gatewaycontract.RequestMetadata{
			RequestId: uuid.NewString(),
		},
		Dto: requestDto,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return exceptions.New(
			"FailedToEncodeRequest",
			"Email",
			"Call",
			"Failed to encode the email service request",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+path,
		bytes.NewReader(payload),
	)
	if err != nil {
		return exceptions.New(
			"FailedToCreateRequest",
			"Email",
			"Call",
			"Failed to create the email service request",
			http.StatusInternalServerError,
			true,
		).WithOrigin(err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return exceptions.New(
			"EmailServiceUnavailable",
			"Email",
			"Call",
			"The email service is unavailable",
			http.StatusServiceUnavailable,
			true,
		).WithOrigin(err)
	}
	defer httpResponse.Body.Close()

	var response gatewaycontract.Response[emaildto.SendEmailResponseDto]
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return exceptions.New(
			"InvalidServiceResponse",
			"Email",
			"Call",
			"The email service returned an invalid response",
			http.StatusBadGateway,
			true,
		).WithOrigin(err)
	}
	if response.Exception != nil {
		return response.Exception
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return exceptions.New(
			"EmailServiceRequestFailed",
			"Email",
			"Call",
			fmt.Sprintf("The email service returned HTTP status %d", httpResponse.StatusCode),
			httpResponse.StatusCode,
			true,
		)
	}

	return nil
}
