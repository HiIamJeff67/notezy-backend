package realtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

func TestGatewaySendsReadyAndPong(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/"+constants.RealtimeDevelopmentBaseURL, NewGateway().Handle)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/" + constants.RealtimeDevelopmentBaseURL
	connection, response, err := websocket.DefaultDialer.Dial(
		wsURL,
		http.Header{"Origin": []string{server.URL}},
	)
	if err != nil {
		t.Fatalf("failed to connect to gateway: %v", err)
	}
	defer connection.Close()
	if response.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("expected status %d, got %d", http.StatusSwitchingProtocols, response.StatusCode)
	}

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

func TestGatewayMultiplexesBlockPackChannels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/"+constants.RealtimeDevelopmentBaseURL, NewGateway().Handle)

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/" + constants.RealtimeDevelopmentBaseURL
	connection, _, err := websocket.DefaultDialer.Dial(
		wsURL,
		http.Header{"Origin": []string{server.URL}},
	)
	if err != nil {
		t.Fatalf("failed to connect to gateway: %v", err)
	}
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
	for _, subscribe := range []realtimetypes.SubscribeFrame{
		{
			Version:     constants.RealtimeProtocolVersion,
			Type:        realtimetypes.FrameType_Subscribe,
			RequestId:   "subscribe-first",
			ChannelType: realtimetypes.ChannelType_BlockPack,
			ChannelId:   firstBlockPackId,
		},
		{
			Version:     constants.RealtimeProtocolVersion,
			Type:        realtimetypes.FrameType_Subscribe,
			RequestId:   "subscribe-second",
			ChannelType: realtimetypes.ChannelType_BlockPack,
			ChannelId:   secondBlockPackId,
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
		Version:     constants.RealtimeProtocolVersion,
		Type:        realtimetypes.FrameType_Subscribe,
		RequestId:   "subscribe-first-again",
		ChannelType: realtimetypes.ChannelType_BlockPack,
		ChannelId:   firstBlockPackId,
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

	if err := connection.WriteJSON(realtimetypes.AcknowledgeFrame{
		Version:            constants.RealtimeProtocolVersion,
		Type:               realtimetypes.FrameType_Acknowledge,
		RequestId:          "ack-second",
		ConnectorChannelId: secondSubscribed.ConnectorChannelId,
		Sequence:           7,
	}); err != nil {
		t.Fatalf("failed to acknowledge channel: %v", err)
	}

	var acknowledged realtimetypes.AcknowledgedFrame
	if err := connection.ReadJSON(&acknowledged); err != nil {
		t.Fatalf("failed to read acknowledged frame: %v", err)
	}
	if acknowledged.ConnectorChannelId != secondSubscribed.ConnectorChannelId || acknowledged.Sequence != 7 {
		t.Fatalf("unexpected acknowledged frame: %#v", acknowledged)
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

	var binaryError realtimetypes.ErrorFrame
	if err := connection.ReadJSON(&binaryError); err != nil {
		t.Fatalf("failed to read binary frame error: %v", err)
	}
	if binaryError.Code != realtimetypes.ErrorCode_BinaryChannelNotReady ||
		binaryError.ConnectorChannelId != firstSubscribed.ConnectorChannelId {
		t.Fatalf("unexpected binary frame error: %#v", binaryError)
	}

	if err := connection.WriteJSON(realtimetypes.UnsubscribeFrame{
		Version:            constants.RealtimeProtocolVersion,
		Type:               realtimetypes.FrameType_Unsubscribe,
		RequestId:          "unsubscribe-second",
		ConnectorChannelId: secondSubscribed.ConnectorChannelId,
	}); err != nil {
		t.Fatalf("failed to unsubscribe channel: %v", err)
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
