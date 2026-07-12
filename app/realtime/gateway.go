package realtime

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	workers "github.com/HiIamJeff67/notezy-backend/app/realtime/workers"
	tokens "github.com/HiIamJeff67/notezy-backend/app/tokens"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Gateway struct {
	upgrader       websocket.Upgrader
	workerManager  workers.WorkerManagerInterface
	connectors     map[uuid.UUID]*Connector
	connectorMutex sync.RWMutex
}

func NewGateway() *Gateway {
	workerManager := workers.NewWorkerManager()
	gateway := &Gateway{
		workerManager: workerManager,
		connectors:    make(map[uuid.UUID]*Connector),
	}
	gateway.upgrader = websocket.Upgrader{
		CheckOrigin: func(req *http.Request) bool {
			return req.Header.Get("Origin") != ""
		},
	}
	workerManager.SetFrameHandler(gateway.handleInternalFrame)

	return gateway
}

func (g *Gateway) Handle(ctx *gin.Context) {
	connectionTicket := websocket.Subprotocols(ctx.Request)
	if len(connectionTicket) != 1 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	connectionClaims, err := tokens.ParseRealtimeConnectionTicket(
		connectionTicket[0],
		ctx.GetHeader("User-Agent"),
	)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userPublicId, err := uuid.Parse(connectionClaims.Subject)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	websocketConnection, err := g.upgrader.Upgrade(
		ctx.Writer,
		ctx.Request,
		http.Header{"Sec-WebSocket-Protocol": []string{connectionTicket[0]}},
	)
	if err != nil {
		return
	}
	defer websocketConnection.Close()

	connector := Connector{
		Id:           uuid.New(),
		UserPublicId: userPublicId,
		UserAgent:    ctx.GetHeader("User-Agent"),
		connection:   websocketConnection,
		channels:     make(map[uint32]realtimetypes.Channel),
	}
	g.connectorMutex.Lock()
	g.connectors[connector.Id] = &connector
	g.connectorMutex.Unlock()

	defer func() {
		g.connectorMutex.Lock()
		delete(g.connectors, connector.Id)
		g.connectorMutex.Unlock()
	}()
	defer func() {
		connector.channelMutex.Lock()
		defer connector.channelMutex.Unlock()
		for connectorChannelId, channel := range connector.channels {
			g.workerManager.Detach(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_Detach,
				ChannelType:        channel.Type,
				ConnectionId:       connector.Id,
				ConnectorChannelId: connectorChannelId,
				ChannelId:          channel.Id,
			})
		}
	}()

	websocketConnection.SetReadLimit(constants.RealtimeMaxMessageSize)
	websocketConnection.SetReadDeadline(time.Now().Add(constants.RealtimePongWait))
	websocketConnection.SetPongHandler(func(string) error {
		return websocketConnection.SetReadDeadline(time.Now().Add(constants.RealtimePongWait))
	})

	if err := connector.writeJSON(realtimetypes.ReadyFrame{
		Version:             constants.RealtimeProtocolVersion,
		Type:                realtimetypes.FrameType_Ready,
		ConnectionId:        connector.Id.String(),
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

	internalFrameType := realtimetypes.InternalFrameType_YjsDocument
	if frame.Type == realtimetypes.BinaryFrameType_Awareness {
		internalFrameType = realtimetypes.InternalFrameType_Awareness
	}

	if !g.workerManager.Forward(realtimetypes.InternalFrame{
		Version:            byte(constants.RealtimeWorkerProtocolVersion),
		Type:               internalFrameType,
		ChannelType:        channel.Type,
		ConnectionId:       connector.Id,
		ConnectorChannelId: frame.ConnectorChannelId,
		ChannelId:          channel.Id,
		Payload:            frame.Payload,
	}) {
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Error,
			ChannelType:        channel.Type,
			ChannelId:          &channel.Id,
			ConnectorChannelId: frame.ConnectorChannelId,
			Code:               realtimetypes.ErrorCode_WorkerUnavailable,
			Message:            "the yjs worker is unavailable",
		})
	}

	return true
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

		channel := realtimetypes.Channel{
			Type: subscribeFrame.ChannelType,
			Id:   subscribeFrame.ChannelId,
		}

		channelClaims, err := tokens.ParseRealtimeBlockPackTicket(
			subscribeFrame.ChannelTicket,
			connector.UserAgent,
		)
		if err != nil || channelClaims.Subject != connector.UserPublicId.String() ||
			channelClaims.ChannelType != string(channel.Type) || channelClaims.ChannelId != channel.Id.String() {
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:     constants.RealtimeProtocolVersion,
				Type:        realtimetypes.FrameType_Error,
				RequestId:   subscribeFrame.RequestId,
				ChannelType: channel.Type,
				ChannelId:   &channel.Id,
				Code:        realtimetypes.ErrorCode_InvalidChannelTicket,
				Message:     "channel ticket is invalid",
			})
		}
		connectorChannelId, existing := connector.subscribe(channel)
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

		if existing {
			return connector.writeJSON(realtimetypes.SubscribedFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Subscribed,
				RequestId:          subscribeFrame.RequestId,
				ChannelType:        subscribeFrame.ChannelType,
				ChannelId:          subscribeFrame.ChannelId,
				ConnectorChannelId: connectorChannelId,
				Existing:           true,
			}) == nil
		}

		if !g.workerManager.Attach(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_Attach,
			ChannelType:        channel.Type,
			ConnectionId:       connector.Id,
			ConnectorChannelId: connectorChannelId,
			ChannelId:          channel.Id,
		}) {
			connector.unsubscribe(connectorChannelId)

			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          subscribeFrame.RequestId,
				ChannelType:        channel.Type,
				ChannelId:          &channel.Id,
				ConnectorChannelId: connectorChannelId,
				Code:               realtimetypes.ErrorCode_WorkerUnavailable,
				Message:            "the yjs worker is unavailable",
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

		channel, exists := connector.unsubscribe(unsubscribeFrame.ConnectorChannelId)

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

		g.workerManager.Detach(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_Detach,
			ChannelType:        channel.Type,
			ConnectionId:       connector.Id,
			ConnectorChannelId: unsubscribeFrame.ConnectorChannelId,
			ChannelId:          channel.Id,
		})

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

func (g *Gateway) handleInternalFrame(frame realtimetypes.InternalFrame) {
	g.connectorMutex.RLock()
	connector, exists := g.connectors[frame.ConnectionId]
	g.connectorMutex.RUnlock()

	if !exists {
		return
	}

	channel, exists := connector.get(frame.ConnectorChannelId)
	if !exists || channel.Type != frame.ChannelType || channel.Id != frame.ChannelId {
		return
	}

	if frame.Type == realtimetypes.InternalFrameType_ResyncRequired ||
		frame.Type == realtimetypes.InternalFrameType_PermissionRevoked {
		connector.unsubscribe(frame.ConnectorChannelId)
		g.workerManager.Detach(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_Detach,
			ChannelType:        channel.Type,
			ConnectionId:       connector.Id,
			ConnectorChannelId: frame.ConnectorChannelId,
			ChannelId:          channel.Id,
		})

		code := realtimetypes.ErrorCode_ResubscribeRequired
		message := "the yjs worker requires this channel to resubscribe"
		if frame.Type == realtimetypes.InternalFrameType_PermissionRevoked {
			code = realtimetypes.ErrorCode_PermissionRevoked
			message = "permission for this channel has been revoked"
		}

		connector.writeError(realtimetypes.ErrorFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Error,
			ChannelType:        channel.Type,
			ChannelId:          &channel.Id,
			ConnectorChannelId: frame.ConnectorChannelId,
			Code:               code,
			Message:            message,
		})

		return
	}

	binaryFrameType := realtimetypes.BinaryFrameType_YjsDocument
	if frame.Type == realtimetypes.InternalFrameType_Awareness {
		binaryFrameType = realtimetypes.BinaryFrameType_Awareness
	} else if frame.Type != realtimetypes.InternalFrameType_YjsDocument {
		return
	}

	payload, err := realtimetypes.BinaryFrame{
		Version:            byte(constants.RealtimeProtocolVersion),
		Type:               binaryFrameType,
		ConnectorChannelId: frame.ConnectorChannelId,
		Payload:            frame.Payload,
	}.MarshalBytes()
	if err != nil {
		return
	}

	connector.writeBinary(payload)
}
