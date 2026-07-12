package realtime

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	tokens "github.com/HiIamJeff67/notezy-backend/app/tokens"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type fakeWorkerManager struct {
	frameHandler func(realtimetypes.InternalFrame)
	frames       []realtimetypes.InternalFrame
	mutex        sync.Mutex
}

type fakeBlockProjectionService struct {
	blockPackId uuid.UUID
	input       dtos.ApplyBlockProjectionInput
}

func (s *fakeBlockProjectionService) Apply(
	ctx context.Context,
	blockPackId uuid.UUID,
	input dtos.ApplyBlockProjectionInput,
) (*dtos.ApplyBlockProjectionResult, error) {
	s.blockPackId = blockPackId
	s.input = input

	return &dtos.ApplyBlockProjectionResult{
		Applied:                true,
		ProjectedUntilSequence: input.ProjectedSequence,
	}, nil
}

func (m *fakeWorkerManager) Attach(frame realtimetypes.InternalFrame) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.frames = append(m.frames, frame)

	return true
}

func (m *fakeWorkerManager) Detach(frame realtimetypes.InternalFrame) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.frames = append(m.frames, frame)
}

func (m *fakeWorkerManager) Forward(frame realtimetypes.InternalFrame) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.frames = append(m.frames, frame)
	if m.frameHandler != nil {
		m.frameHandler(frame)
	}

	return true
}

func (m *fakeWorkerManager) SetFrameHandler(handler func(realtimetypes.InternalFrame)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.frameHandler = handler
}

func TestGatewaySendsReadyAndPong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userAgent := "notezy-realtime-test"
	userPublicId := uuid.New()
	configureRealtimeTicketPrivateKey(t)
	connectionTicket, _, exception := tokens.GenerateRealtimeConnectionTicket(userPublicId, userAgent)
	if exception != nil {
		t.Fatalf("failed to generate connection ticket: %v", exception)
	}

	workerManager := &fakeWorkerManager{}
	gateway := &Gateway{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(req *http.Request) bool {
				return req.Header.Get("Origin") != ""
			},
		},
		workerManager: workerManager,
		connectors:    make(map[uuid.UUID]*Connector),
	}
	workerManager.SetFrameHandler(gateway.handleInternalFrame)

	router := gin.New()
	router.GET("/"+constants.RealtimeDevelopmentBaseURL, gateway.Handle)

	server := httptest.NewServer(router)
	defer server.Close()

	connection := dialGateway(t, server.URL, userAgent, *connectionTicket)
	defer connection.Close()

	var ready realtimetypes.ReadyFrame
	if err := connection.ReadJSON(&ready); err != nil {
		t.Fatalf("failed to read ready frame: %v", err)
	}
	if ready.Version != constants.RealtimeProtocolVersion ||
		ready.Type != realtimetypes.FrameType_Ready ||
		ready.ConnectionId == "" {
		t.Fatalf("unexpected ready frame: %#v", ready)
	}

	if err := connection.WriteJSON(realtimetypes.ControlFrame{
		Version:   constants.RealtimeProtocolVersion,
		Type:      realtimetypes.FrameType_Ping,
		RequestId: "request-1",
	}); err != nil {
		t.Fatalf("failed to write ping frame: %v", err)
	}

	var pong realtimetypes.ControlFrame
	if err := connection.ReadJSON(&pong); err != nil {
		t.Fatalf("failed to read pong frame: %v", err)
	}
	if pong.Version != constants.RealtimeProtocolVersion ||
		pong.Type != realtimetypes.FrameType_Pong ||
		pong.RequestId != "request-1" {
		t.Fatalf("unexpected pong frame: %#v", pong)
	}
}

func TestGatewayMultiplexesAndRelaysBlockPackChannels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userAgent := "notezy-realtime-test"
	userPublicId := uuid.New()
	configureRealtimeTicketPrivateKey(t)
	connectionTicket, _, exception := tokens.GenerateRealtimeConnectionTicket(userPublicId, userAgent)
	if exception != nil {
		t.Fatalf("failed to generate connection ticket: %v", exception)
	}

	workerManager := &fakeWorkerManager{}
	gateway := &Gateway{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(req *http.Request) bool {
				return req.Header.Get("Origin") != ""
			},
		},
		workerManager: workerManager,
		connectors:    make(map[uuid.UUID]*Connector),
	}
	workerManager.SetFrameHandler(gateway.handleInternalFrame)

	router := gin.New()
	router.GET("/"+constants.RealtimeDevelopmentBaseURL, gateway.Handle)

	server := httptest.NewServer(router)
	defer server.Close()

	connection := dialGateway(t, server.URL, userAgent, *connectionTicket)
	defer connection.Close()

	var ready realtimetypes.ReadyFrame
	if err := connection.ReadJSON(&ready); err != nil {
		t.Fatalf("failed to read ready frame: %v", err)
	}
	if !ready.ResubscribeRequired {
		t.Fatalf("expected ready frame to require resubscription: %#v", ready)
	}

	unsupportedChannelId := uuid.New()
	if err := connection.WriteJSON(realtimetypes.SubscribeFrame{
		Version:     constants.RealtimeProtocolVersion,
		Type:        realtimetypes.FrameType_Subscribe,
		RequestId:   "subscribe-unsupported",
		ChannelType: realtimetypes.ChannelType("Unsupported"),
		ChannelId:   unsupportedChannelId,
	}); err != nil {
		t.Fatalf("failed to subscribe to unsupported channel type: %v", err)
	}

	var unsupportedChannelTypeError realtimetypes.ErrorFrame
	if err := connection.ReadJSON(&unsupportedChannelTypeError); err != nil {
		t.Fatalf("failed to read unsupported channel type error: %v", err)
	}
	if unsupportedChannelTypeError.Code != realtimetypes.ErrorCode_UnsupportedChannelType ||
		unsupportedChannelTypeError.ChannelId == nil ||
		*unsupportedChannelTypeError.ChannelId != unsupportedChannelId {
		t.Fatalf("unexpected unsupported channel type error: %#v", unsupportedChannelTypeError)
	}

	firstBlockPackId := uuid.New()
	secondBlockPackId := uuid.New()
	channelTickets := make(map[uuid.UUID]string, 2)
	for _, blockPackId := range []uuid.UUID{firstBlockPackId, secondBlockPackId} {
		channelTicket, _, exception := tokens.GenerateRealtimeBlockPackTicket(
			userPublicId,
			userAgent,
			blockPackId,
			realtimetypes.ChannelPermission_Write,
		)
		if exception != nil {
			t.Fatalf("failed to generate channel ticket: %v", exception)
		}

		channelTickets[blockPackId] = *channelTicket
	}

	for _, subscribe := range []realtimetypes.SubscribeFrame{
		{
			Version:       constants.RealtimeProtocolVersion,
			Type:          realtimetypes.FrameType_Subscribe,
			RequestId:     "subscribe-first",
			ChannelType:   realtimetypes.ChannelType_BlockPack,
			ChannelId:     firstBlockPackId,
			ChannelTicket: channelTickets[firstBlockPackId],
		},
		{
			Version:       constants.RealtimeProtocolVersion,
			Type:          realtimetypes.FrameType_Subscribe,
			RequestId:     "subscribe-second",
			ChannelType:   realtimetypes.ChannelType_BlockPack,
			ChannelId:     secondBlockPackId,
			ChannelTicket: channelTickets[secondBlockPackId],
		},
	} {
		if err := connection.WriteJSON(subscribe); err != nil {
			t.Fatalf("failed to subscribe: %v", err)
		}
	}

	var firstSubscribed realtimetypes.SubscribedFrame
	if err := connection.ReadJSON(&firstSubscribed); err != nil {
		t.Fatalf("failed to read first subscribed frame: %v", err)
	}
	var secondSubscribed realtimetypes.SubscribedFrame
	if err := connection.ReadJSON(&secondSubscribed); err != nil {
		t.Fatalf("failed to read second subscribed frame: %v", err)
	}
	if firstSubscribed.ConnectorChannelId == 0 || secondSubscribed.ConnectorChannelId == 0 ||
		firstSubscribed.ConnectorChannelId == secondSubscribed.ConnectorChannelId ||
		firstSubscribed.ChannelType != realtimetypes.ChannelType_BlockPack ||
		secondSubscribed.ChannelType != realtimetypes.ChannelType_BlockPack ||
		firstSubscribed.ChannelId != firstBlockPackId ||
		secondSubscribed.ChannelId != secondBlockPackId {
		t.Fatalf("unexpected subscribed frames: %#v %#v", firstSubscribed, secondSubscribed)
	}

	if err := connection.WriteJSON(realtimetypes.SubscribeFrame{
		Version:       constants.RealtimeProtocolVersion,
		Type:          realtimetypes.FrameType_Subscribe,
		RequestId:     "subscribe-first-again",
		ChannelType:   realtimetypes.ChannelType_BlockPack,
		ChannelId:     firstBlockPackId,
		ChannelTicket: channelTickets[firstBlockPackId],
	}); err != nil {
		t.Fatalf("failed to repeat subscribe: %v", err)
	}

	var existingSubscribed realtimetypes.SubscribedFrame
	if err := connection.ReadJSON(&existingSubscribed); err != nil {
		t.Fatalf("failed to read existing subscribed frame: %v", err)
	}
	if !existingSubscribed.Existing || existingSubscribed.ConnectorChannelId != firstSubscribed.ConnectorChannelId {
		t.Fatalf("expected idempotent subscription: %#v", existingSubscribed)
	}

	binaryPayload, err := realtimetypes.BinaryFrame{
		Version:            byte(constants.RealtimeProtocolVersion),
		Type:               realtimetypes.BinaryFrameType_YjsDocument,
		ConnectorChannelId: firstSubscribed.ConnectorChannelId,
		Payload:            []byte{1, 2, 3},
	}.MarshalBytes()
	if err != nil {
		t.Fatalf("failed to marshal binary frame: %v", err)
	}

	if err := connection.WriteMessage(websocket.BinaryMessage, binaryPayload); err != nil {
		t.Fatalf("failed to write binary frame: %v", err)
	}

	messageType, relayedPayload, err := connection.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read relayed binary frame: %v", err)
	}
	if messageType != websocket.BinaryMessage {
		t.Fatalf("expected relayed binary frame, got message type %d", messageType)
	}

	var relayedFrame realtimetypes.BinaryFrame
	if err := relayedFrame.UnmarshalBytes(relayedPayload); err != nil {
		t.Fatalf("failed to unmarshal relayed binary frame: %v", err)
	}
	if relayedFrame.Type != realtimetypes.BinaryFrameType_YjsDocument ||
		relayedFrame.ConnectorChannelId != firstSubscribed.ConnectorChannelId ||
		string(relayedFrame.Payload) != string([]byte{1, 2, 3}) {
		t.Fatalf("unexpected relayed binary frame: %#v", relayedFrame)
	}

	if err := connection.WriteJSON(realtimetypes.UnsubscribeFrame{
		Version:            constants.RealtimeProtocolVersion,
		Type:               realtimetypes.FrameType_Unsubscribe,
		RequestId:          "unsubscribe-second",
		ConnectorChannelId: secondSubscribed.ConnectorChannelId,
	}); err != nil {
		t.Fatalf("failed to unsubscribe: %v", err)
	}

	var unsubscribed realtimetypes.UnsubscribedFrame
	if err := connection.ReadJSON(&unsubscribed); err != nil {
		t.Fatalf("failed to read unsubscribed frame: %v", err)
	}
	if unsubscribed.ConnectorChannelId != secondSubscribed.ConnectorChannelId ||
		unsubscribed.ChannelType != realtimetypes.ChannelType_BlockPack ||
		unsubscribed.ChannelId != secondBlockPackId {
		t.Fatalf("unexpected unsubscribed frame: %#v", unsubscribed)
	}
}

func TestGatewayRejectsYjsDocumentUpdatesOnReadOnlyChannels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userAgent := "notezy-realtime-test"
	userPublicId := uuid.New()
	configureRealtimeTicketPrivateKey(t)
	connectionTicket, _, exception := tokens.GenerateRealtimeConnectionTicket(userPublicId, userAgent)
	if exception != nil {
		t.Fatalf("failed to generate connection ticket: %v", exception)
	}

	workerManager := &fakeWorkerManager{}
	gateway := &Gateway{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(req *http.Request) bool {
				return req.Header.Get("Origin") != ""
			},
		},
		workerManager: workerManager,
		connectors:    make(map[uuid.UUID]*Connector),
	}
	workerManager.SetFrameHandler(gateway.handleInternalFrame)

	router := gin.New()
	router.GET("/"+constants.RealtimeDevelopmentBaseURL, gateway.Handle)

	server := httptest.NewServer(router)
	defer server.Close()

	connection := dialGateway(t, server.URL, userAgent, *connectionTicket)
	defer connection.Close()

	var ready realtimetypes.ReadyFrame
	if err := connection.ReadJSON(&ready); err != nil {
		t.Fatalf("failed to read ready frame: %v", err)
	}

	blockPackId := uuid.New()
	channelTicket, _, exception := tokens.GenerateRealtimeBlockPackTicket(
		userPublicId,
		userAgent,
		blockPackId,
		realtimetypes.ChannelPermission_Read,
	)
	if exception != nil {
		t.Fatalf("failed to generate read channel ticket: %v", exception)
	}

	if err := connection.WriteJSON(realtimetypes.SubscribeFrame{
		Version:       constants.RealtimeProtocolVersion,
		Type:          realtimetypes.FrameType_Subscribe,
		RequestId:     "subscribe-read",
		ChannelType:   realtimetypes.ChannelType_BlockPack,
		ChannelId:     blockPackId,
		ChannelTicket: *channelTicket,
	}); err != nil {
		t.Fatalf("failed to subscribe to read channel: %v", err)
	}

	var subscribed realtimetypes.SubscribedFrame
	if err := connection.ReadJSON(&subscribed); err != nil {
		t.Fatalf("failed to read subscribed frame: %v", err)
	}

	binaryPayload, err := realtimetypes.BinaryFrame{
		Version:            byte(constants.RealtimeProtocolVersion),
		Type:               realtimetypes.BinaryFrameType_YjsDocument,
		ConnectorChannelId: subscribed.ConnectorChannelId,
		Payload:            []byte{1, 2, 3},
	}.MarshalBytes()
	if err != nil {
		t.Fatalf("failed to marshal Yjs document frame: %v", err)
	}

	if err := connection.WriteMessage(websocket.BinaryMessage, binaryPayload); err != nil {
		t.Fatalf("failed to write Yjs document frame: %v", err)
	}

	var permissionError realtimetypes.ErrorFrame
	if err := connection.ReadJSON(&permissionError); err != nil {
		t.Fatalf("failed to read channel permission error: %v", err)
	}
	if permissionError.Code != realtimetypes.ErrorCode_ChannelPermissionDenied ||
		permissionError.ConnectorChannelId != subscribed.ConnectorChannelId ||
		permissionError.ChannelId == nil || *permissionError.ChannelId != blockPackId {
		t.Fatalf("unexpected channel permission error: %#v", permissionError)
	}

	workerManager.mutex.Lock()
	defer workerManager.mutex.Unlock()
	for _, frame := range workerManager.frames {
		if frame.Type == realtimetypes.InternalFrameType_YjsDocument {
			t.Fatalf("read-only Yjs document frame must not reach the worker: %#v", frame)
		}
	}
}

func TestGatewayAppliesBlockProjectionInternalFrames(t *testing.T) {
	workerManager := &fakeWorkerManager{}
	blockProjectionService := &fakeBlockProjectionService{}
	gateway := &Gateway{
		workerManager:          workerManager,
		blockProjectionService: blockProjectionService,
		connectors:             make(map[uuid.UUID]*Connector),
	}
	workerManager.SetFrameHandler(gateway.handleInternalFrame)

	blockPackId := uuid.New()
	connectionId := uuid.New()
	input := dtos.ApplyBlockProjectionInput{
		SchemaId:          "notezy.blocknote",
		SchemaVersion:     1,
		ProjectedSequence: 7,
		Blocks:            []dtos.ArborizedEditableBlock{},
	}
	payload, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal block projection input: %v", err)
	}

	gateway.handleInternalFrame(realtimetypes.InternalFrame{
		Version:            byte(constants.RealtimeWorkerProtocolVersion),
		Type:               realtimetypes.InternalFrameType_ApplyBlockProjection,
		ChannelType:        realtimetypes.ChannelType_BlockPack,
		ConnectionId:       connectionId,
		ConnectorChannelId: 1,
		ChannelId:          blockPackId,
		Payload:            payload,
	})

	if blockProjectionService.blockPackId != blockPackId ||
		blockProjectionService.input.ProjectedSequence != input.ProjectedSequence {
		t.Fatalf("unexpected projection service invocation: %#v", blockProjectionService)
	}

	workerManager.mutex.Lock()
	defer workerManager.mutex.Unlock()
	if len(workerManager.frames) != 1 {
		t.Fatalf("expected one projection response frame, got %d", len(workerManager.frames))
	}

	frame := workerManager.frames[0]
	if frame.Type != realtimetypes.InternalFrameType_BlockProjectionApplied ||
		frame.ChannelId != blockPackId || frame.ConnectionId != connectionId {
		t.Fatalf("unexpected projection response frame: %#v", frame)
	}

	var result dtos.ApplyBlockProjectionResult
	if err := json.Unmarshal(frame.Payload, &result); err != nil {
		t.Fatalf("failed to unmarshal block projection response: %v", err)
	}
	if !result.Applied || result.ProjectedUntilSequence != input.ProjectedSequence {
		t.Fatalf("unexpected block projection response: %#v", result)
	}
}

func dialGateway(t *testing.T, serverURL string, userAgent string, connectionTicket string) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(serverURL, "http") + "/" + constants.RealtimeDevelopmentBaseURL
	connection, response, err := (&websocket.Dialer{
		Subprotocols: []string{connectionTicket},
	}).Dial(wsURL, http.Header{
		"Origin":     []string{serverURL},
		"User-Agent": []string{userAgent},
	})
	if err != nil {
		t.Fatalf("failed to connect to gateway: %v", err)
	}
	if response.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected status %d, got %d", http.StatusSwitchingProtocols, response.StatusCode)
	}

	return connection
}

func configureRealtimeTicketPrivateKey(t *testing.T) {
	t.Helper()

	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate realtime ticket private key: %v", err)
	}
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("failed to marshal realtime ticket private key: %v", err)
	}

	t.Setenv("REALTIME_TICKET_PRIVATE_KEY_BASE64", base64.StdEncoding.EncodeToString(privateKeyBytes))
}
