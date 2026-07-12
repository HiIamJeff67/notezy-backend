package realtime

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Gateway struct {
	upgrader websocket.Upgrader
}

func NewGateway() *Gateway {
	return &Gateway{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(req *http.Request) bool {
				return req.Header.Get("Origin") != ""
			},
		},
	}
}

func (g *Gateway) Handle(ctx *gin.Context) {
	websocketConnection, err := g.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer websocketConnection.Close()

	connector := Connector{
		connection: websocketConnection,
		channels:   make(map[uint32]Channel),
	}

	websocketConnection.SetReadLimit(constants.RealtimeMaxMessageSize)
	websocketConnection.SetReadDeadline(time.Now().Add(constants.RealtimePongWait))
	websocketConnection.SetPongHandler(func(string) error {
		return websocketConnection.SetReadDeadline(time.Now().Add(constants.RealtimePongWait))
	})

	if err := connector.writeJSON(realtimetypes.ReadyFrame{
		Version:             constants.RealtimeProtocolVersion,
		Type:                realtimetypes.FrameType_Ready,
		ConnectionId:        uuid.NewString(),
		ResubscribeRequired: true,
	}); err != nil {
		return
	}

	pingDone := make(chan struct{})
	defer close(pingDone)

	go func() {
		ticker := time.NewTicker(constants.RealtimePingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-pingDone:
				return
			case <-ticker.C:
				if err := connector.writeControl(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	for {
		messageType, payload, err := websocketConnection.ReadMessage()
		if err != nil {
			return
		}

		switch messageType {
		case websocket.BinaryMessage:
			if !g.handleBinaryFrame(&connector, payload) {
				return
			}
		case websocket.TextMessage:
			if !g.handleControlFrame(&connector, payload) {
				return
			}
		default:
			if !connector.writeError(realtimetypes.ErrorFrame{
				Version: constants.RealtimeProtocolVersion,
				Type:    realtimetypes.FrameType_Error,
				Code:    realtimetypes.ErrorCode_UnsupportedMessageType,
				Message: "only text control frames and binary channel frames are supported",
			}) {
				return
			}
		}
	}
}

func (g *Gateway) handleBinaryFrame(connector *Connector, payload []byte) bool {
	var frame realtimetypes.BinaryFrame
	if err := frame.UnmarshalBytes(payload); err != nil {
		return connector.writeError(realtimetypes.ErrorFrame{
			Version: constants.RealtimeProtocolVersion,
			Type:    realtimetypes.FrameType_Error,
			Code:    realtimetypes.ErrorCode_InvalidBinaryFrame,
			Message: "binary frames must contain a version, type, channelId, and payload",
		})
	}
	if int(frame.Version) != constants.RealtimeProtocolVersion {
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Error,
			ConnectorChannelId: frame.ConnectorChannelId,
			Code:               realtimetypes.ErrorCode_UnsupportedProtocolVersion,
			Message:            "unsupported realtime protocol version",
		})
	}

	channel, exists := connector.get(frame.ConnectorChannelId)

	if !exists {
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Error,
			ConnectorChannelId: frame.ConnectorChannelId,
			Code:               realtimetypes.ErrorCode_ChannelNotFound,
			Message:            "connectorChannelId is not subscribed on this connection",
		})
	}
	if frame.Type != realtimetypes.BinaryFrameType_YjsDocument &&
		frame.Type != realtimetypes.BinaryFrameType_Awareness {
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Error,
			ChannelType:        channel.Type,
			ChannelId:          &channel.Id,
			ConnectorChannelId: frame.ConnectorChannelId,
			Code:               realtimetypes.ErrorCode_UnsupportedBinaryType,
			Message:            "binary frame type is not enabled",
		})
	}

	return connector.writeError(realtimetypes.ErrorFrame{
		Version:            constants.RealtimeProtocolVersion,
		Type:               realtimetypes.FrameType_Error,
		ChannelType:        channel.Type,
		ChannelId:          &channel.Id,
		ConnectorChannelId: frame.ConnectorChannelId,
		Code:               realtimetypes.ErrorCode_BinaryChannelNotReady,
		Message:            "the yjs worker channel is not enabled yet",
	})
}

func (g *Gateway) handleControlFrame(connector *Connector, payload []byte) bool {
	var controlFrame realtimetypes.ControlFrame
	if err := json.Unmarshal(payload, &controlFrame); err != nil {
		return connector.writeError(realtimetypes.ErrorFrame{
			Version: constants.RealtimeProtocolVersion,
			Type:    realtimetypes.FrameType_Error,
			Code:    realtimetypes.ErrorCode_InvalidControlFrame,
			Message: "control frames must be valid JSON",
		})
	}
	if controlFrame.Version != constants.RealtimeProtocolVersion {
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:   constants.RealtimeProtocolVersion,
			Type:      realtimetypes.FrameType_Error,
			RequestId: controlFrame.RequestId,
			Code:      realtimetypes.ErrorCode_UnsupportedProtocolVersion,
			Message:   "unsupported realtime protocol version",
		})
	}

	switch controlFrame.Type {
	case realtimetypes.FrameType_Ping:
		return connector.writeJSON(realtimetypes.ControlFrame{
			Version:   constants.RealtimeProtocolVersion,
			Type:      realtimetypes.FrameType_Pong,
			RequestId: controlFrame.RequestId,
		}) == nil
	case realtimetypes.FrameType_Heartbeat:
		return connector.writeJSON(realtimetypes.HeartbeatFrame{
			Version:      constants.RealtimeProtocolVersion,
			Type:         realtimetypes.FrameType_Heartbeat,
			RequestId:    controlFrame.RequestId,
			UnixMilliNow: time.Now().UnixMilli(),
		}) == nil
	case realtimetypes.FrameType_Subscribe:
		var subscribeFrame realtimetypes.SubscribeFrame
		if err := json.Unmarshal(payload, &subscribeFrame); err != nil || subscribeFrame.ChannelId == uuid.Nil {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:   constants.RealtimeProtocolVersion,
				Type:      realtimetypes.FrameType_Error,
				RequestId: controlFrame.RequestId,
				Code:      realtimetypes.ErrorCode_InvalidChannelId,
				Message:   "subscribe requires a valid channelId",
			})
		}

		switch subscribeFrame.ChannelType {
		case realtimetypes.ChannelType_BlockPack:
		default:
			if subscribeFrame.ChannelType == "" {
				return connector.writeError(realtimetypes.ErrorFrame{
					Version:   constants.RealtimeProtocolVersion,
					Type:      realtimetypes.FrameType_Error,
					RequestId: controlFrame.RequestId,
					ChannelId: &subscribeFrame.ChannelId,
					Code:      realtimetypes.ErrorCode_InvalidChannelType,
					Message:   "subscribe requires a channelType",
				})
			}

			return connector.writeError(realtimetypes.ErrorFrame{
				Version:     constants.RealtimeProtocolVersion,
				Type:        realtimetypes.FrameType_Error,
				RequestId:   controlFrame.RequestId,
				ChannelType: subscribeFrame.ChannelType,
				ChannelId:   &subscribeFrame.ChannelId,
				Code:        realtimetypes.ErrorCode_UnsupportedChannelType,
				Message:     "channelType is not enabled",
			})
		}

		channel := Channel{
			Type: subscribeFrame.ChannelType,
			Id:   subscribeFrame.ChannelId,
		}
		connectorChannelId, existing := connector.append(channel)
		if connectorChannelId == 0 {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:     constants.RealtimeProtocolVersion,
				Type:        realtimetypes.FrameType_Error,
				RequestId:   subscribeFrame.RequestId,
				ChannelType: subscribeFrame.ChannelType,
				ChannelId:   &subscribeFrame.ChannelId,
				Code:        realtimetypes.ErrorCode_ChannelLimitExceeded,
				Message:     "the connection cannot subscribe to more channels",
			})
		}

		return connector.writeJSON(realtimetypes.SubscribedFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Subscribed,
			RequestId:          subscribeFrame.RequestId,
			ChannelType:        subscribeFrame.ChannelType,
			ChannelId:          subscribeFrame.ChannelId,
			ConnectorChannelId: connectorChannelId,
			Existing:           existing,
		}) == nil
	case realtimetypes.FrameType_Unsubscribe:
		var unsubscribeFrame realtimetypes.UnsubscribeFrame
		if err := json.Unmarshal(payload, &unsubscribeFrame); err != nil || unsubscribeFrame.ConnectorChannelId == 0 {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:   constants.RealtimeProtocolVersion,
				Type:      realtimetypes.FrameType_Error,
				RequestId: controlFrame.RequestId,
				Code:      realtimetypes.ErrorCode_InvalidConnectorChannelId,
				Message:   "unsubscribe requires a valid connectorChannelId",
			})
		}

		channel, exists := connector.remove(unsubscribeFrame.ConnectorChannelId)

		if !exists {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          unsubscribeFrame.RequestId,
				ConnectorChannelId: unsubscribeFrame.ConnectorChannelId,
				Code:               realtimetypes.ErrorCode_ChannelNotFound,
				Message:            "connectorChannelId is not subscribed on this connection",
			})
		}

		return connector.writeJSON(realtimetypes.UnsubscribedFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Unsubscribed,
			RequestId:          unsubscribeFrame.RequestId,
			ChannelType:        channel.Type,
			ChannelId:          channel.Id,
			ConnectorChannelId: unsubscribeFrame.ConnectorChannelId,
		}) == nil
	case realtimetypes.FrameType_Acknowledge:
		var acknowledgeFrame realtimetypes.AcknowledgeFrame
		if err := json.Unmarshal(payload, &acknowledgeFrame); err != nil || acknowledgeFrame.ConnectorChannelId == 0 {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:   constants.RealtimeProtocolVersion,
				Type:      realtimetypes.FrameType_Error,
				RequestId: controlFrame.RequestId,
				Code:      realtimetypes.ErrorCode_InvalidConnectorChannelId,
				Message:   "ack requires a valid connectorChannelId",
			})
		}

		exists, acknowledged := connector.acknowledge(
			acknowledgeFrame.ConnectorChannelId,
			acknowledgeFrame.Sequence,
		)

		if !exists {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          acknowledgeFrame.RequestId,
				ConnectorChannelId: acknowledgeFrame.ConnectorChannelId,
				Code:               realtimetypes.ErrorCode_ChannelNotFound,
				Message:            "connectorChannelId is not subscribed on this connection",
			})
		}
		if !acknowledged {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          acknowledgeFrame.RequestId,
				ConnectorChannelId: acknowledgeFrame.ConnectorChannelId,
				Code:               realtimetypes.ErrorCode_InvalidAcknowledgement,
				Message:            "ack sequence cannot move backwards",
			})
		}

		return connector.writeJSON(realtimetypes.AcknowledgedFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Acknowledged,
			RequestId:          acknowledgeFrame.RequestId,
			ConnectorChannelId: acknowledgeFrame.ConnectorChannelId,
			Sequence:           acknowledgeFrame.Sequence,
		}) == nil
	case realtimetypes.FrameType_Reconnect:
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:   constants.RealtimeProtocolVersion,
			Type:      realtimetypes.FrameType_Error,
			RequestId: controlFrame.RequestId,
			Code:      realtimetypes.ErrorCode_ResubscribeRequired,
			Message:   "new connections must resubscribe their channels",
		})
	case realtimetypes.FrameType_Authenticate:
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:   constants.RealtimeProtocolVersion,
			Type:      realtimetypes.FrameType_Error,
			RequestId: controlFrame.RequestId,
			Code:      realtimetypes.ErrorCode_AuthenticationManagedByUpgrade,
			Message:   "authenticate with the WebSocket upgrade request",
		})
	default:
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:   constants.RealtimeProtocolVersion,
			Type:      realtimetypes.FrameType_Error,
			RequestId: controlFrame.RequestId,
			Code:      realtimetypes.ErrorCode_UnsupportedControlType,
			Message:   "control frame type is not enabled",
		})
	}
}
