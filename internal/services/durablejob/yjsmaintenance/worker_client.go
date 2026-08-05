package yjsmaintenance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"
	yjsworkercontract "github.com/HiIamJeff67/notezy-backend/contracts/yjsworker/v1"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/traces"
	durablejobconfig "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/config"
)

type WorkerClient struct {
	documentInitializationEndpoint string
	client                         *http.Client
}

func NewWorkerClient(config durablejobconfig.Config) WorkerClient {
	return WorkerClient{
		documentInitializationEndpoint: config.YjsDocumentInitializationWorkerUrl,
		client: &http.Client{
			Timeout: yjsWorkerRequestTimeout,
		},
	}
}

func (c WorkerClient) InitializeDocuments(
	ctx context.Context,
	reqDtos []blockpacksdto.InitializeBlockPackYjsDocumentReqDto,
) (resDtos []blockpacksdto.InitializeBlockPackYjsDocumentResDto, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.worker.initialize_documents")
	span.SetAttributes(attribute.Int("yjs.document_count", len(reqDtos)))
	defer func() {
		outcome := "success"
		if err != nil {
			outcome = "error"
			logs.NotezyLogger.Error(ctx, err, "Yjs document initialization worker request failed", attribute.String("operation", "worker.initialize_documents"))
		}
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "worker.initialize_documents"),
			attribute.String("outcome", outcome),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "worker.initialize_documents"),
			attribute.String("outcome", outcome),
		)
		traces.NotezyTracer.End(span, err)
	}()

	if len(reqDtos) == 0 {
		return []blockpacksdto.InitializeBlockPackYjsDocumentResDto{}, nil
	}
	if c.documentInitializationEndpoint == "" {
		return nil, fmt.Errorf("YJS_DOCUMENT_INITIALIZATION_WORKER_URL is required")
	}

	payload, err := json.Marshal(struct {
		Documents []blockpacksdto.InitializeBlockPackYjsDocumentReqDto `json:"documents"`
	}{
		Documents: reqDtos,
	})
	if err != nil {
		return nil, err
	}
	if len(payload) > yjsworkercontract.YjsMaintenanceMaximumPayloadBytes {
		return nil, fmt.Errorf("yjs document initialization batch exceeds the worker payload limit")
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.documentInitializationEndpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yjs document initialization worker returned %s", response.Status)
	}

	responsePayload, err := io.ReadAll(io.LimitReader(response.Body, int64(yjsworkercontract.YjsMaintenanceMaximumPayloadBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(responsePayload) > yjsworkercontract.YjsMaintenanceMaximumPayloadBytes {
		return nil, fmt.Errorf("yjs document initialization worker response exceeds the payload limit")
	}

	var responseDto struct {
		Documents []blockpacksdto.InitializeBlockPackYjsDocumentResDto `json:"documents"`
	}
	if err := json.Unmarshal(responsePayload, &responseDto); err != nil {
		return nil, err
	}
	if len(responseDto.Documents) != len(reqDtos) {
		return nil, errors.New("incomplete yjs document initialization worker response")
	}
	for _, document := range responseDto.Documents {
		if len(document.Snapshot) == 0 || len(document.StateVector) == 0 {
			return nil, errors.New("invalid yjs document initialization worker response")
		}
	}

	return responseDto.Documents, nil
}
