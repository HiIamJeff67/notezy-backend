package yjsmaintenance

import (
	"bytes"
	"context"
	"encoding/binary"
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

	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type WorkerClient struct {
	compactionEndpoint string
	projectionEndpoint string
	client             *http.Client
}

func NewWorkerClient() WorkerClient {
	return WorkerClient{
		compactionEndpoint: os.Getenv("YJS_MAINTENANCE_WORKER_URL"),
		projectionEndpoint: os.Getenv("YJS_PROJECTION_WORKER_URL"),
		client: &http.Client{
			Timeout: constants.YjsMaintenanceWorkerRequestTimeout,
		},
	}
}

func (c WorkerClient) Compact(
	ctx context.Context,
	inputs []realtimetypes.YjsCompactionBatchInput,
) (results []realtimetypes.YjsCompactionBatchResult, err error) {
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

	results = make([]realtimetypes.YjsCompactionBatchResult, 0, resultCount)
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

		var result realtimetypes.YjsCompactionBatchResult
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
	inputs []realtimetypes.YjsProjectionBatchInput,
) (results []realtimetypes.YjsProjectionBatchResult, err error) {
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

	results = make([]realtimetypes.YjsProjectionBatchResult, 0, resultCount)
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

		var result realtimetypes.YjsProjectionBatchResult
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
