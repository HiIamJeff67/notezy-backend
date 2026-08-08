package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gatewaycontract "github.com/HiIamJeff67/notezy-backend/contracts/gateway/v1"
)

func TestCoreAdapterForwardsVersionedEnvelopeAndMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Authorization") != "Bearer delegation-token" {
			t.Fatal("expected delegation token header")
		}
		if request.Header.Get("Traceparent") != "00-trace" {
			t.Fatal("expected trace parent header")
		}
		if request.Header.Get("Cookie") != "" {
			t.Fatal("cookies must not cross the Gateway/Core boundary")
		}

		requestEnvelope := gatewaycontract.Request[struct{}]{}
		if err := json.NewDecoder(request.Body).Decode(&requestEnvelope); err != nil {
			t.Fatalf("decode request envelope: %v", err)
		}
		if requestEnvelope.Version != gatewaycontract.Version || requestEnvelope.Operation != "station.get" {
			t.Fatal("expected versioned station request envelope")
		}
		if requestEnvelope.Tokens.AccessToken != "access-token" {
			t.Fatalf("expected typed access token, got %q", requestEnvelope.Tokens.AccessToken)
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(responseWriter).Encode(&gatewaycontract.Response[struct{}]{
			Version: gatewaycontract.Version,
			Metadata: gatewaycontract.ResponseMetadata{
				RequestId: requestEnvelope.Metadata.RequestId,
			},
			Data: struct{}{},
		}); err != nil {
			t.Fatalf("encode response envelope: %v", err)
		}
	}))
	defer server.Close()

	client := NewCoreAdapter(server.URL, time.Second)
	response, err := call[struct{}, struct{}](
		client,
		nil,
		context.Background(),
		http.MethodPost,
		"/v1/operations",
		"delegation-token",
		http.Header{
			"Cookie":     []string{"accessToken=test-token"},
			"User-Agent": []string{"test-agent"},
			"X-Real-IP":  []string{"192.0.2.1"},
		},
		&gatewaycontract.Request[struct{}]{
			Operation: "station.get",
			Metadata: gatewaycontract.RequestMetadata{
				RequestId:   "request-id",
				TraceParent: "00-trace",
			},
			Tokens: gatewaycontract.Tokens{
				AccessToken: "access-token",
			},
		},
	)
	if err != nil {
		t.Fatalf("execute Core service request: %v", err)
	}
	if response.Metadata.RequestId != "request-id" {
		t.Fatalf("expected request ID request-id, got %s", response.Metadata.RequestId)
	}
}
