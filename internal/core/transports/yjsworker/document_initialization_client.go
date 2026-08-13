package adapterstransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	blockpackscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"
	yjsworkercontract "github.com/HiIamJeff67/notezy-backend/contracts/yjs-worker/v1"

	coreconfig "github.com/HiIamJeff67/notezy-backend/internal/core/configs"
)

type DocumentInitializationClient struct {
	endpoint   string
	httpClient *http.Client
}

func NewDocumentInitializationClient(
	config coreconfig.YjsDocumentInitializationConfig,
) *DocumentInitializationClient {
	return &DocumentInitializationClient{
		endpoint: config.Endpoint,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *DocumentInitializationClient) InitializeDocuments(
	ctx context.Context,
	requestDtos []blockpackscontract.InitializeBlockPackYjsDocumentReqDto,
) ([]blockpackscontract.InitializeBlockPackYjsDocumentResDto, error) {
	if len(requestDtos) == 0 {
		return []blockpackscontract.InitializeBlockPackYjsDocumentResDto{}, nil
	}

	payload, err := json.Marshal(struct {
		Documents []blockpackscontract.InitializeBlockPackYjsDocumentReqDto `json:"documents"`
	}{
		Documents: requestDtos,
	})
	if err != nil {
		return nil, fmt.Errorf("encode Yjs document initialization request: %w", err)
	}
	if len(payload) > yjsworkercontract.YjsMaintenanceMaximumPayloadBytes {
		return nil, errors.New("Yjs document initialization request exceeds the worker payload limit")
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("create Yjs document initialization request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send Yjs document initialization request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Yjs document initialization worker returned %s", response.Status)
	}

	responsePayload, err := io.ReadAll(io.LimitReader(
		response.Body,
		int64(yjsworkercontract.YjsMaintenanceMaximumPayloadBytes)+1,
	))
	if err != nil {
		return nil, fmt.Errorf("read Yjs document initialization response: %w", err)
	}
	if len(responsePayload) > yjsworkercontract.YjsMaintenanceMaximumPayloadBytes {
		return nil, errors.New("Yjs document initialization response exceeds the worker payload limit")
	}

	var responseDto struct {
		Documents []blockpackscontract.InitializeBlockPackYjsDocumentResDto `json:"documents"`
	}
	if err := json.Unmarshal(responsePayload, &responseDto); err != nil {
		return nil, fmt.Errorf("decode Yjs document initialization response: %w", err)
	}
	if len(responseDto.Documents) != len(requestDtos) {
		return nil, errors.New("Yjs document initialization response is incomplete")
	}
	for _, document := range responseDto.Documents {
		if len(document.Snapshot) == 0 || len(document.StateVector) == 0 {
			return nil, errors.New("Yjs document initialization response is invalid")
		}
	}

	return responseDto.Documents, nil
}
