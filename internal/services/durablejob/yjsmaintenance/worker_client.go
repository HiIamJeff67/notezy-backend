package yjsmaintenance

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

	blockpacksdto "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/api/block-packs"
	yjsworkercontract "github.com/HiIamJeff67/notezy-backend/contracts/yjsworker/v1"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/traces"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type WorkerClient struct {
	compactionEndpoint             string
	documentInitializationEndpoint string
	projectionEndpoint             string
	client                         *http.Client
}

func NewWorkerClient() WorkerClient {
	return WorkerClient{
		compactionEndpoint:             os.Getenv("YJS_MAINTENANCE_WORKER_URL"),
		documentInitializationEndpoint: os.Getenv("YJS_DOCUMENT_INITIALIZATION_WORKER_URL"),
		projectionEndpoint:             os.Getenv("YJS_PROJECTION_WORKER_URL"),
		client: &http.Client{
			Timeout: constants.YjsMaintenanceWorkerRequestTimeout,
		},
	}
}

func (c WorkerClient) InitializeDocuments(
	ctx context.Context,
	reqDtos []blockpacksdto.InitializeBlockPackYjsDocumentReqDto,
) (resDtos []blockpacksdto.InitializeBlockPackYjsDocumentResDto, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.maintenance.worker.initialize_documents")
	span.SetAttributes(attribute.Int("yjs.document_count", len(reqDtos)))
	defer func() {
		outcome := "success"
		if err != nil {
			outcome = "error"
			logs.NotezyLogger.Error(ctx, err, "Yjs document initialization worker request failed", attribute.String("operation", "maintenance.worker.initialize_documents"))
		}
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "maintenance.worker.initialize_documents"),
			attribute.String("outcome", outcome),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "maintenance.worker.initialize_documents"),
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
	if len(payload) > constants.YjsMaintenanceWorkerMaxPayloadBytes {
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

	responsePayload, err := io.ReadAll(io.LimitReader(response.Body, int64(constants.YjsMaintenanceWorkerMaxPayloadBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(responsePayload) > constants.YjsMaintenanceWorkerMaxPayloadBytes {
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

func (c WorkerClient) Compact(
	ctx context.Context,
	inputs []yjsworkercontract.YjsCompactionBatchInput,
) (results []yjsworkercontract.YjsCompactionBatchResult, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.maintenance.worker.compact")
	span.SetAttributes(attribute.Int("yjs.document_count", len(inputs)))
	defer func() {
		outcome := "success"
		if err != nil {
			outcome = "error"
			logs.NotezyLogger.Error(ctx, err, "Yjs maintenance worker request failed", attribute.String("operation", "maintenance.worker.compact"))
		}
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "maintenance.worker.compact"),
			attribute.String("outcome", outcome),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "maintenance.worker.compact"),
			attribute.String("outcome", outcome),
		)
		traces.NotezyTracer.End(span, err)
	}()

	if len(inputs) == 0 {
		return nil, nil
	}
	if c.compactionEndpoint == "" {
		return nil, fmt.Errorf("YJS_MAINTENANCE_WORKER_URL is required")
	}

	payload := bytes.NewBuffer(make([]byte, 0))
	if err := binary.Write(payload, binary.BigEndian, uint32(len(inputs))); err != nil {
		return nil, err
	}

	blockPackIdSet := make(map[[16]byte]bool, len(inputs))
	for _, input := range inputs {
		if blockPackIdSet[input.BlockPackId] {
			return nil, errors.New("duplicate yjs maintenance block pack id")
		}
		blockPackIdSet[input.BlockPackId] = true

		inputPayload, err := input.MarshalBytes()
		if err != nil {
			return nil, err
		}
		if len(inputPayload) > math.MaxUint32 || payload.Len()+4+len(inputPayload) > constants.YjsMaintenanceWorkerMaxPayloadBytes {
			return nil, fmt.Errorf("yjs maintenance batch exceeds the worker payload limit")
		}

		if err := binary.Write(payload, binary.BigEndian, uint32(len(inputPayload))); err != nil {
			return nil, err
		}
		if _, err := payload.Write(inputPayload); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.compactionEndpoint,
		bytes.NewReader(payload.Bytes()),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/octet-stream")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yjs maintenance worker returned %s", response.Status)
	}

	responsePayload, err := io.ReadAll(io.LimitReader(response.Body, int64(constants.YjsMaintenanceWorkerMaxPayloadBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(responsePayload) > constants.YjsMaintenanceWorkerMaxPayloadBytes {
		return nil, fmt.Errorf("yjs maintenance worker response exceeds the payload limit")
	}

	if len(responsePayload) < 4 {
		return nil, errors.New("invalid yjs maintenance worker response")
	}

	resultCount := binary.BigEndian.Uint32(responsePayload[0:4])
	if resultCount != uint32(len(inputs)) {
		return nil, errors.New("incomplete yjs maintenance worker response")
	}

	results = make([]yjsworkercontract.YjsCompactionBatchResult, 0, resultCount)
	resultBlockPackIdSet := make(map[[16]byte]bool, resultCount)
	offset := 4
	for index := uint32(0); index < resultCount; index++ {
		if len(responsePayload)-offset < 4 {
			return nil, errors.New("invalid yjs maintenance worker response")
		}

		resultLength := binary.BigEndian.Uint32(responsePayload[offset : offset+4])
		offset += 4
		if uint64(resultLength) > uint64(len(responsePayload)-offset) {
			return nil, errors.New("invalid yjs maintenance worker response")
		}

		var result yjsworkercontract.YjsCompactionBatchResult
		if err := result.UnmarshalBytes(responsePayload[offset : offset+int(resultLength)]); err != nil {
			return nil, err
		}
		if resultBlockPackIdSet[result.BlockPackId] {
			return nil, errors.New("duplicate yjs maintenance worker result")
		}
		resultBlockPackIdSet[result.BlockPackId] = true
		offset += int(resultLength)

		results = append(results, result)
	}
	if offset != len(responsePayload) {
		return nil, errors.New("invalid yjs maintenance worker response")
	}

	return results, nil
}

func (c WorkerClient) Project(
	ctx context.Context,
	inputs []yjsworkercontract.YjsProjectionBatchInput,
) (results []yjsworkercontract.YjsProjectionBatchResult, err error) {
	start := time.Now()
	ctx, span := traces.NotezyTracer.Start(ctx, "yjs.maintenance.worker.project")
	span.SetAttributes(attribute.Int("yjs.document_count", len(inputs)))
	defer func() {
		outcome := "success"
		if err != nil {
			outcome = "error"
			logs.NotezyLogger.Error(ctx, err, "Yjs projection worker request failed", attribute.String("operation", "maintenance.worker.project"))
		}
		metrics.NotezyMeter.Count(ctx, "yjs.operation.count", 1,
			attribute.String("operation", "maintenance.worker.project"),
			attribute.String("outcome", outcome),
		)
		metrics.NotezyMeter.Duration(ctx, "yjs.operation.duration", time.Since(start),
			attribute.String("operation", "maintenance.worker.project"),
			attribute.String("outcome", outcome),
		)
		traces.NotezyTracer.End(span, err)
	}()

	if len(inputs) == 0 {
		return nil, nil
	}
	if c.projectionEndpoint == "" {
		return nil, fmt.Errorf("YJS_PROJECTION_WORKER_URL is required")
	}

	payload := bytes.NewBuffer(make([]byte, 0))
	if err := binary.Write(payload, binary.BigEndian, uint32(len(inputs))); err != nil {
		return nil, err
	}
	inputBlockPackIdSet := make(map[[16]byte]bool, len(inputs))
	for _, input := range inputs {
		if inputBlockPackIdSet[input.BlockPackId] {
			return nil, errors.New("duplicate yjs projection block pack id")
		}
		inputBlockPackIdSet[input.BlockPackId] = true

		inputPayload, err := input.MarshalBytes()
		if err != nil {
			return nil, err
		}
		if len(inputPayload) > math.MaxUint32 || payload.Len()+4+len(inputPayload) > constants.YjsMaintenanceWorkerMaxPayloadBytes {
			return nil, fmt.Errorf("yjs projection batch exceeds the worker payload limit")
		}

		if err := binary.Write(payload, binary.BigEndian, uint32(len(inputPayload))); err != nil {
			return nil, err
		}
		if _, err := payload.Write(inputPayload); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.projectionEndpoint, bytes.NewReader(payload.Bytes()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/octet-stream")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(request.Header))

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yjs projection worker returned %s", response.Status)
	}

	responsePayload, err := io.ReadAll(io.LimitReader(response.Body, int64(constants.YjsMaintenanceWorkerMaxPayloadBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(responsePayload) > constants.YjsMaintenanceWorkerMaxPayloadBytes || len(responsePayload) < 4 {
		return nil, errors.New("invalid yjs projection worker response")
	}

	resultCount := binary.BigEndian.Uint32(responsePayload[:4])
	if resultCount != uint32(len(inputs)) {
		return nil, errors.New("incomplete yjs projection worker response")
	}

	results = make([]yjsworkercontract.YjsProjectionBatchResult, 0, resultCount)
	resultBlockPackIdSet := make(map[[16]byte]bool, resultCount)
	offset := 4
	for index := uint32(0); index < resultCount; index++ {
		if len(responsePayload)-offset < 4 {
			return nil, errors.New("invalid yjs projection worker response")
		}

		resultLength := binary.BigEndian.Uint32(responsePayload[offset : offset+4])
		offset += 4
		if uint64(resultLength) > uint64(len(responsePayload)-offset) {
			return nil, errors.New("invalid yjs projection worker response")
		}

		var result yjsworkercontract.YjsProjectionBatchResult
		if err := result.UnmarshalBytes(responsePayload[offset : offset+int(resultLength)]); err != nil {
			return nil, err
		}
		if resultBlockPackIdSet[result.BlockPackId] || !inputBlockPackIdSet[result.BlockPackId] {
			return nil, errors.New("invalid yjs projection worker result")
		}
		resultBlockPackIdSet[result.BlockPackId] = true
		offset += int(resultLength)
		results = append(results, result)
	}
	if offset != len(responsePayload) || len(resultBlockPackIdSet) != len(inputBlockPackIdSet) {
		return nil, errors.New("invalid yjs projection worker response")
	}

	return results, nil
}
