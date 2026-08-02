package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	corecontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1"
	durablejobdto "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1"
	blocksdto "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1/api/blocks"
	sharedtokens "github.com/HiIamJeff67/notezy-backend/internal/shared/tokens"
)

const applyBlockProjectionPath = "/durablejob/" + durablejobdto.ApplyBlockProjectionOperation

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() Client {
	baseURL := strings.TrimRight(os.Getenv("CORE_BASE_URL"), "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:7778"
	}

	return Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c Client) ApplyBlockProjections(
	ctx context.Context,
	documents []blocksdto.ApplyBlockProjectionDocumentRequestDto,
) ([]uuid.UUID, error) {
	requestId := uuid.NewString()
	token, err := sharedtokens.GenerateDelegationToken(sharedtokens.DelegationTokenClaims{
		Actor:     "durablejob",
		Operation: durablejobdto.ApplyBlockProjectionOperation,
		RequestId: requestId,
	})
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(corecontract.Request[durablejobdto.ApplyBlockProjectionRequestDto]{
		Version:   corecontract.Version,
		Operation: durablejobdto.ApplyBlockProjectionOperation,
		Metadata: corecontract.RequestMetadata{
			RequestId: requestId,
		},
		Dto: durablejobdto.ApplyBlockProjectionRequestDto{
			Documents: documents,
		},
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+applyBlockProjectionPath,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+*token)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result corecontract.Response[durablejobdto.ApplyBlockProjectionResponseDto]
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Exception != nil {
		return nil, fmt.Errorf("core projection failed: %s", result.Exception.String())
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("core projection returned %s", response.Status)
	}

	return result.Data.AppliedBlockPackIds, nil
}
