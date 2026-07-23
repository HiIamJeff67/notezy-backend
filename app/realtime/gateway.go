package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"

	caches "github.com/HiIamJeff67/notezy-backend/app/caches"
	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	models "github.com/HiIamJeff67/notezy-backend/app/models"
	repositories "github.com/HiIamJeff67/notezy-backend/app/models/repositories"
	scopes "github.com/HiIamJeff67/notezy-backend/app/models/scopes"
	logs "github.com/HiIamJeff67/notezy-backend/app/monitor/logs"
	metrics "github.com/HiIamJeff67/notezy-backend/app/monitor/metrics"
	traces "github.com/HiIamJeff67/notezy-backend/app/monitor/traces"
	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	workers "github.com/HiIamJeff67/notezy-backend/app/realtime/workers"
	services "github.com/HiIamJeff67/notezy-backend/app/services"
	tokens "github.com/HiIamJeff67/notezy-backend/app/tokens"
	util "github.com/HiIamJeff67/notezy-backend/app/util"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Gateway struct {
	upgrader              websocket.Upgrader
	workerManager         workers.WorkerManagerInterface
	yjsPersistenceService services.YjsPersistenceServiceInterface
	realtimeService       interface {
		GetBlockPackChannelPermission(ctx context.Context, userPublicId uuid.UUID, blockPackId uuid.UUID, permission realtimetypes.ChannelPermission) (int32, realtimetypes.ErrorCode, error)
		ValidateBlockPackChannelPermission(ctx context.Context, userPublicId uuid.UUID, blockPackId uuid.UUID, permission realtimetypes.ChannelPermission) (realtimetypes.ErrorCode, error)
	}
	leaseStore             *caches.RealtimeLeaseStore
	blockProjectionService interface {
		Apply(ctx context.Context, blockPackId uuid.UUID, input dtos.ApplyBlockProjectionInput) (*dtos.ApplyBlockProjectionResult, error)
	}
	realtimeDisabled            bool
	realtimeBetaUserPublicIdSet map[uuid.UUID]bool
	connectorMutex              sync.RWMutex
	connectors                  map[uuid.UUID]*Connector
	pendingConnectorCount       int
	maximumConnectors           int
	maximumConnectionsPerUser   int
}

func NewGateway() *Gateway {
	workerManager := workers.NewWorkerManager()
	blockScope := scopes.NewBlockScope()
	blockPackScope := scopes.NewBlockPackScope()
	subShelfScope := scopes.NewSubShelfScope()
	blockPackRepository := repositories.NewBlockPackRepository(blockPackScope)

	realtimeEnabled, err := strconv.ParseBool(util.GetEnv("REALTIME_ENABLED", "true"))
	if err != nil {
		realtimeEnabled = true
	}

	var realtimeBetaUserPublicIdSet map[uuid.UUID]bool
	if rawUserPublicIds := strings.TrimSpace(util.GetEnv("REALTIME_BETA_USER_PUBLIC_IDS", "")); rawUserPublicIds != "" {
		realtimeBetaUserPublicIdSet = make(map[uuid.UUID]bool)
		for _, rawUserPublicId := range strings.Split(rawUserPublicIds, ",") {
			userPublicId, err := uuid.Parse(strings.TrimSpace(rawUserPublicId))
			if err == nil {
				realtimeBetaUserPublicIdSet[userPublicId] = true
			}
		}
	}

	gateway := &Gateway{
		workerManager:         workerManager,
		yjsPersistenceService: services.NewYjsPersistenceService(models.NotezyDB),
		realtimeService:       services.NewRealtimeService(models.NotezyDB, blockPackRepository),
		leaseStore:            caches.NewRealtimeLeaseStore(caches.RedisClientMap),
		blockProjectionService: services.NewBlockService(
			models.NotezyDB,
			blockScope,
			blockPackScope,
			subShelfScope,
			blockPackRepository,
			repositories.NewBlockRepository(blockScope),
		),
		realtimeDisabled:            !realtimeEnabled,
		realtimeBetaUserPublicIdSet: realtimeBetaUserPublicIdSet,
		connectors:                  make(map[uuid.UUID]*Connector),
		maximumConnectors:           constants.RealtimeMaxConnectorsPerGateway,
		maximumConnectionsPerUser:   constants.RealtimeMaxConnectionsPerUser,
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
	// extract and validate the ticket which is in Sec-WebSocket-Protocol header
	connectionTicket := websocket.Subprotocols(ctx.Request)
	if len(connectionTicket) != 1 {
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "missing_connection_ticket"),
		)
		logs.NotezyLogger.Warn(ctx.Request.Context(), fmt.Sprintf("Rejected realtime connection: expected one connection ticket, got %d subprotocols", len(connectionTicket)))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	connectionClaims, err := tokens.ParseRealtimeConnectionTicket(
		connectionTicket[0],
		ctx.GetHeader("User-Agent"),
	)
	if err != nil {
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "invalid_connection_ticket"),
		)
		logs.NotezyLogger.Warn(ctx.Request.Context(), fmt.Sprintf("Rejected realtime connection: invalid connection ticket: %v", err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userPublicId, err := uuid.Parse(connectionClaims.Subject)
	if err != nil {
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "invalid_user_public_id"),
		)
		logs.NotezyLogger.Warn(ctx.Request.Context(), fmt.Sprintf("Rejected realtime connection: connection ticket subject is not a user public id: %v", err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if g.realtimeDisabled ||
		len(g.realtimeBetaUserPublicIdSet) > 0 && !g.realtimeBetaUserPublicIdSet[userPublicId] {
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "rollout_not_enabled"),
		)
		logs.NotezyLogger.Info(ctx.Request.Context(), "Rejected realtime connection because rollout is not enabled")
		ctx.AbortWithStatus(http.StatusServiceUnavailable)

		return
	}

	connectorId := uuid.New()

	maximumConnectors := g.maximumConnectors
	if maximumConnectors <= 0 {
		maximumConnectors = constants.RealtimeMaxConnectorsPerGateway
	}

	g.connectorMutex.Lock()
	if len(g.connectors)+g.pendingConnectorCount >= maximumConnectors {
		g.connectorMutex.Unlock()
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "gateway_capacity_exceeded"),
		)
		ctx.AbortWithStatus(http.StatusServiceUnavailable)

		return
	}
	g.pendingConnectorCount++
	g.connectorMutex.Unlock()

	pendingConnectorAdmission := true
	defer func() {
		if !pendingConnectorAdmission {
			return
		}

		g.connectorMutex.Lock()
		g.pendingConnectorCount--
		g.connectorMutex.Unlock()
	}()

	maximumConnectionsPerUser := g.maximumConnectionsPerUser
	if maximumConnectionsPerUser <= 0 {
		maximumConnectionsPerUser = constants.RealtimeMaxConnectionsPerUser
	}

	acquired, _, err := g.leaseStore.AcquireUserConnection(
		userPublicId,
		connectorId,
		maximumConnectionsPerUser,
	)
	if err != nil {
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "admission_unavailable"),
		)
		logs.NotezyLogger.Error(ctx.Request.Context(), err, "Failed to acquire realtime user connection lease")
		ctx.AbortWithStatus(http.StatusServiceUnavailable)

		return
	}
	if !acquired {
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "user_connection_limit_exceeded"),
		)
		ctx.Header("Retry-After", strconv.Itoa(int(constants.RealtimeLeaseTTL.Seconds())))
		ctx.AbortWithStatus(http.StatusTooManyRequests)

		return
	}
	defer func() {
		if err := g.leaseStore.ReleaseUserConnection(userPublicId, connectorId); err != nil {
			logs.NotezyLogger.Error(ctx.Request.Context(), err, "Failed to release realtime user connection lease")
		}
	}()

	websocketConnection, err := g.upgrader.Upgrade(
		ctx.Writer,
		ctx.Request,
		http.Header{"Sec-WebSocket-Protocol": []string{connectionTicket[0]}},
	)
	if err != nil {
		metrics.NotezyMeter.Count(ctx.Request.Context(), "realtime.connection.rejected.count", 1,
			attribute.String("reason", "websocket_upgrade_failed"),
		)
		traces.NotezyTracer.RecordError(ctx.Request.Context(), err)

		return
	}
	defer websocketConnection.Close()

	connectionStart := time.Now()
	connectionContext, connectionSpan := traces.NotezyTracer.Start(
		ctx.Request.Context(), "realtime.connection",
	)
	defer func() { traces.NotezyTracer.End(connectionSpan, nil) }()

	connector := Connector{
		Id:           connectorId,
		UserPublicId: userPublicId,
		UserAgent:    ctx.GetHeader("User-Agent"),
		connection:   websocketConnection,
		channels:     make(map[uint32]realtimetypes.Channel),
		outbound:     newOutboundQueue(websocketConnection),
	}
	connectionSpan.SetAttributes(attribute.String("realtime.connection.id", connector.Id.String()))
	connector.startWriter()
	defer connector.stopWriter()

	g.connectorMutex.Lock()
	g.pendingConnectorCount--
	g.connectors[connector.Id] = &connector
	g.connectorMutex.Unlock()
	pendingConnectorAdmission = false
	metrics.NotezyMeter.Count(connectionContext, "realtime.connection.accepted.count", 1)
	metrics.NotezyMeter.UpDown(connectionContext, "realtime.connector.count", 1)

	defer func() {
		g.connectorMutex.Lock()
		delete(g.connectors, connector.Id)
		g.connectorMutex.Unlock()
		metrics.NotezyMeter.Count(connectionContext, "realtime.connection.closed.count", 1)
		metrics.NotezyMeter.Duration(connectionContext, "realtime.connection.duration", time.Since(connectionStart))
		metrics.NotezyMeter.UpDown(connectionContext, "realtime.connector.count", -1)
	}()
	defer func() {
		connector.channelMutex.Lock()
		channels := connector.channels
		connector.channels = make(map[uint32]realtimetypes.Channel)
		connector.channelMutex.Unlock()

		for connectorChannelId, channel := range channels {
			if err := g.leaseStore.ReleaseBlockPackSubscriber(
				channel.Id,
				fmt.Sprintf("%s:%d", connector.Id, connectorChannelId),
			); err != nil {
				logs.NotezyLogger.Error(connectionContext, err, "Failed to release realtime BlockPack subscriber lease")
			}

			g.workerManager.Detach(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_Detach,
				ChannelType:        channel.Type,
				ConnectionId:       connector.Id,
				ConnectorChannelId: connectorChannelId,
				ChannelId:          channel.Id,
			})
			metrics.NotezyMeter.Count(connectionContext, "realtime.channel.subscription.count", 1,
				attribute.String("action", "detach"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", "connection_closed"),
			)
			metrics.NotezyMeter.UpDown(connectionContext, "realtime.channel.count", -1,
				attribute.String("channelType", string(channel.Type)),
				attribute.String("permission", string(channel.Permission)),
			)
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
				refreshed, err := g.leaseStore.RefreshUserConnection(connector.UserPublicId, connector.Id)
				if err != nil || !refreshed {
					if err != nil {
						logs.NotezyLogger.Error(connectionContext, err, "Failed to refresh realtime user connection lease")
					}
					_ = websocketConnection.Close()

					return
				}

				connector.channelMutex.RLock()
				channels := make(map[uint32]realtimetypes.Channel, len(connector.channels))
				for connectorChannelId, channel := range connector.channels {
					channels[connectorChannelId] = channel
				}
				connector.channelMutex.RUnlock()

				for connectorChannelId, channel := range channels {
					refreshed, err := g.leaseStore.RefreshBlockPackSubscriber(
						channel.Id,
						fmt.Sprintf("%s:%d", connector.Id, connectorChannelId),
					)
					if err != nil || !refreshed {
						if err != nil {
							logs.NotezyLogger.Error(connectionContext, err, "Failed to refresh realtime BlockPack subscriber lease")
						}
						_ = websocketConnection.Close()

						return
					}
				}

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
			if !g.handleBinaryFrame(connectionContext, &connector, payload) {
				return
			}
		case websocket.TextMessage:
			if !g.handleControlFrame(connectionContext, &connector, payload) {
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

func (g *Gateway) handleBinaryFrame(ctx context.Context, connector *Connector, payload []byte) bool {
	var frame realtimetypes.BinaryFrame
	if err := frame.UnmarshalBytes(payload); err != nil {
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "invalid_binary_frame"),
		)
		return connector.writeError(realtimetypes.ErrorFrame{
			Version: constants.RealtimeProtocolVersion,
			Type:    realtimetypes.FrameType_Error,
			Code:    realtimetypes.ErrorCode_InvalidBinaryFrame,
			Message: "binary frames must contain a version, type, channelId, and payload",
		})
	}
	if int(frame.Version) != constants.RealtimeProtocolVersion {
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "unsupported_protocol_version"),
		)
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
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "channel_not_found"),
		)
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
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "unsupported_binary_type"),
			attribute.String("channelType", string(channel.Type)),
		)
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
	if frame.Type == realtimetypes.BinaryFrameType_YjsDocument &&
		channel.Permission != realtimetypes.ChannelPermission_Write {
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "permission_denied"),
			attribute.String("channelType", string(channel.Type)),
		)
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:            constants.RealtimeProtocolVersion,
			Type:               realtimetypes.FrameType_Error,
			ChannelType:        channel.Type,
			ChannelId:          &channel.Id,
			ConnectorChannelId: frame.ConnectorChannelId,
			Code:               realtimetypes.ErrorCode_ChannelPermissionDenied,
			Message:            "channel permission does not allow yjs document updates",
		})
	}
	if frame.Type == realtimetypes.BinaryFrameType_YjsDocument {
		errorCode, err := g.realtimeService.ValidateBlockPackChannelPermission(
			ctx,
			connector.UserPublicId,
			channel.Id,
			channel.Permission,
		)
		if err != nil {
			connector.unsubscribe(frame.ConnectorChannelId)
			if errorCode == "" {
				errorCode = realtimetypes.ErrorCode_PermissionRevoked
			}

			if releaseErr := g.leaseStore.ReleaseBlockPackSubscriber(
				channel.Id,
				fmt.Sprintf("%s:%d", connector.Id, frame.ConnectorChannelId),
			); releaseErr != nil {
				logs.NotezyLogger.Error(ctx, releaseErr, "Failed to release realtime BlockPack subscriber lease")
			}

			g.workerManager.Detach(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_Detach,
				ChannelType:        channel.Type,
				ConnectionId:       connector.Id,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          channel.Id,
			})

			outcome := "permission_revoked"
			message := "permission for this channel has been revoked"
			if errorCode == realtimetypes.ErrorCode_ResourceUnavailable {
				outcome = "resource_unavailable"
				message = "the block pack is no longer available"
			}
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "detach"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", outcome),
			)
			metrics.NotezyMeter.UpDown(ctx, "realtime.channel.count", -1,
				attribute.String("channelType", string(channel.Type)),
				attribute.String("permission", string(channel.Permission)),
			)

			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				ChannelType:        channel.Type,
				ChannelId:          &channel.Id,
				ConnectorChannelId: frame.ConnectorChannelId,
				Code:               errorCode,
				Message:            message,
			})
		}
	}

	internalFrameType := realtimetypes.InternalFrameType_YjsDocument
	if frame.Type == realtimetypes.BinaryFrameType_Awareness {
		internalFrameType = realtimetypes.InternalFrameType_Awareness
	}
	metrics.NotezyMeter.Count(ctx, "realtime.frame.count", 1,
		attribute.String("direction", "inbound"),
		attribute.String("channelType", string(channel.Type)),
		attribute.String("frameType", string(frame.Type)),
	)
	metrics.NotezyMeter.Bytes(ctx, "realtime.payload.bytes", int64(len(frame.Payload)),
		attribute.String("direction", "inbound"),
		attribute.String("channelType", string(channel.Type)),
		attribute.String("frameType", string(frame.Type)),
	)

	if !g.workerManager.Forward(realtimetypes.InternalFrame{
		Version:            byte(constants.RealtimeWorkerProtocolVersion),
		Type:               internalFrameType,
		ChannelType:        channel.Type,
		ConnectionId:       connector.Id,
		ConnectorChannelId: frame.ConnectorChannelId,
		ChannelId:          channel.Id,
		Payload:            frame.Payload,
	}) {
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "worker_unavailable"),
			attribute.String("channelType", string(channel.Type)),
		)
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

func (g *Gateway) handleControlFrame(ctx context.Context, connector *Connector, payload []byte) bool {
	var controlFrame realtimetypes.ControlFrame
	if err := json.Unmarshal(payload, &controlFrame); err != nil {
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "invalid_control_frame"),
		)
		return connector.writeError(realtimetypes.ErrorFrame{
			Version: constants.RealtimeProtocolVersion,
			Type:    realtimetypes.FrameType_Error,
			Code:    realtimetypes.ErrorCode_InvalidControlFrame,
			Message: "control frames must be valid JSON",
		})
	}
	if controlFrame.Version != constants.RealtimeProtocolVersion {
		metrics.NotezyMeter.Count(ctx, "realtime.frame.rejected.count", 1,
			attribute.String("direction", "inbound"),
			attribute.String("reason", "unsupported_protocol_version"),
		)
		return connector.writeError(realtimetypes.ErrorFrame{
			Version:   constants.RealtimeProtocolVersion,
			Type:      realtimetypes.FrameType_Error,
			RequestId: controlFrame.RequestId,
			Code:      realtimetypes.ErrorCode_UnsupportedProtocolVersion,
			Message:   "unsupported realtime protocol version",
		})
	}
	metrics.NotezyMeter.Count(ctx, "realtime.frame.count", 1,
		attribute.String("direction", "inbound"),
		attribute.String("frameType", string(controlFrame.Type)),
	)
	metrics.NotezyMeter.Bytes(ctx, "realtime.payload.bytes", int64(len(payload)),
		attribute.String("direction", "inbound"),
		attribute.String("frameType", string(controlFrame.Type)),
	)

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
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("outcome", "invalid_channel_id"),
			)
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
				metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
					attribute.String("action", "subscribe"),
					attribute.String("outcome", "invalid_channel_type"),
				)
				return connector.writeError(realtimetypes.ErrorFrame{
					Version:   constants.RealtimeProtocolVersion,
					Type:      realtimetypes.FrameType_Error,
					RequestId: controlFrame.RequestId,
					ChannelId: &subscribeFrame.ChannelId,
					Code:      realtimetypes.ErrorCode_InvalidChannelType,
					Message:   "subscribe requires a channelType",
				})
			}
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(subscribeFrame.ChannelType)),
				attribute.String("outcome", "unsupported_channel_type"),
			)

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

		channelClaims, err := tokens.ParseRealtimeBlockPackTicket(
			subscribeFrame.ChannelTicket,
			connector.UserAgent,
		)
		if err != nil || channelClaims.Subject != connector.UserPublicId.String() ||
			channelClaims.ChannelType != string(subscribeFrame.ChannelType) ||
			channelClaims.ChannelId != subscribeFrame.ChannelId.String() {
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(subscribeFrame.ChannelType)),
				attribute.String("outcome", "invalid_channel_ticket"),
			)
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:     constants.RealtimeProtocolVersion,
				Type:        realtimetypes.FrameType_Error,
				RequestId:   subscribeFrame.RequestId,
				ChannelType: subscribeFrame.ChannelType,
				ChannelId:   &subscribeFrame.ChannelId,
				Code:        realtimetypes.ErrorCode_InvalidChannelTicket,
				Message:     "channel ticket is invalid",
			})
		}

		// create the channel here, so if handleControlFrame of subscribe does not fire first
		// the channel just will not be found by g.connectors.get() methods, and error will be thrown
		channel := realtimetypes.Channel{
			Type:       subscribeFrame.ChannelType,
			Id:         subscribeFrame.ChannelId,
			Permission: realtimetypes.ChannelPermission(channelClaims.Permission),
		}
		connectorChannelId, existing := connector.subscribe(channel)
		if connectorChannelId == 0 {
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", "channel_limit_exceeded"),
			)
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
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", "existing"),
			)
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

		// else if the channel is not exist yet, we have to create a channel by first check the maximum subscribers
		maximumSubscribers, errorCode, err := g.realtimeService.GetBlockPackChannelPermission(
			ctx,
			connector.UserPublicId,
			channel.Id,
			channel.Permission,
		)
		if err != nil {
			connector.unsubscribe(connectorChannelId)
			if errorCode == "" {
				errorCode = realtimetypes.ErrorCode_PermissionRevoked
			}

			message := "permission for this channel has been revoked"
			if errorCode == realtimetypes.ErrorCode_ResourceUnavailable {
				message = "the block pack is no longer available"
			} else if errorCode == realtimetypes.ErrorCode_RoomAdmissionUnavailable {
				message = "room admission is temporarily unavailable"
			}

			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", string(errorCode)),
			)

			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          subscribeFrame.RequestId,
				ChannelType:        channel.Type,
				ChannelId:          &channel.Id,
				ConnectorChannelId: connectorChannelId,
				Code:               errorCode,
				Message:            message,
			})
		}

		leaseMember := fmt.Sprintf("%s:%d", connector.Id, connectorChannelId)
		acquired, activeSubscribers, err := g.leaseStore.AcquireBlockPackSubscriber(
			channel.Id,
			leaseMember,
			int(maximumSubscribers),
		)
		if err != nil {
			connector.unsubscribe(connectorChannelId)
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", "admission_unavailable"),
			)
			logs.NotezyLogger.Error(ctx, err, "Failed to acquire realtime BlockPack subscriber lease")

			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          subscribeFrame.RequestId,
				ChannelType:        channel.Type,
				ChannelId:          &channel.Id,
				ConnectorChannelId: connectorChannelId,
				Code:               realtimetypes.ErrorCode_RoomAdmissionUnavailable,
				Message:            "room admission is temporarily unavailable",
			})
		}
		if !acquired {
			connector.unsubscribe(connectorChannelId)

			leaseMembers := make([]string, 0)
			leases, err := g.leaseStore.GetBlockPackSubscriberLeases(channel.Id)
			if err != nil {
				logs.NotezyLogger.Error(ctx, err, "Failed to inspect realtime BlockPack subscriber leases")
			} else {
				for _, lease := range leases {
					leaseMembers = append(leaseMembers, fmt.Sprintf("%s expiresAt=%s", lease.Member, lease.ExpiresAt.UTC().Format(time.RFC3339Nano)))
				}
			}
			logs.NotezyLogger.Warn(ctx, "Rejected realtime BlockPack subscription because subscriber limit was reached",
				attribute.String("realtime.block_pack.id", channel.Id.String()),
				attribute.Int("realtime.room.maximum_subscribers", int(maximumSubscribers)),
				attribute.Int64("realtime.room.active_subscribers", activeSubscribers),
				attribute.StringSlice("realtime.room.subscriber_leases", leaseMembers),
			)

			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", "room_connection_limit_exceeded"),
			)

			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          subscribeFrame.RequestId,
				ChannelType:        channel.Type,
				ChannelId:          &channel.Id,
				ConnectorChannelId: connectorChannelId,
				Code:               realtimetypes.ErrorCode_RoomConnectionLimitExceeded,
				Message:            "the room has reached the active subscriber limit for its plan",
			})
		}
		if err := g.leaseStore.SetBlockPackParticipant(
			channel.Id,
			leaseMember,
			connector.UserPublicId,
			string(channel.Permission),
		); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to record realtime BlockPack participant")
		}

		if !g.workerManager.Attach(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_Attach,
			ChannelType:        channel.Type,
			ConnectionId:       connector.Id,
			ConnectorChannelId: connectorChannelId,
			ChannelId:          channel.Id,
		}) {
			if err := g.leaseStore.ReleaseBlockPackSubscriber(channel.Id, leaseMember); err != nil {
				logs.NotezyLogger.Error(ctx, err, "Failed to release realtime BlockPack subscriber lease")
			}

			connector.unsubscribe(connectorChannelId)
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "subscribe"),
				attribute.String("channelType", string(channel.Type)),
				attribute.String("outcome", "worker_unavailable"),
			)

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
		metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
			attribute.String("action", "subscribe"),
			attribute.String("channelType", string(channel.Type)),
			attribute.String("outcome", "success"),
		)
		metrics.NotezyMeter.UpDown(ctx, "realtime.channel.count", 1,
			attribute.String("channelType", string(channel.Type)),
			attribute.String("permission", string(channel.Permission)),
		)

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
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "unsubscribe"),
				attribute.String("outcome", "invalid_connector_channel_id"),
			)
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
			metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
				attribute.String("action", "unsubscribe"),
				attribute.String("outcome", "channel_not_found"),
			)
			return connector.writeError(realtimetypes.ErrorFrame{
				Version:            constants.RealtimeProtocolVersion,
				Type:               realtimetypes.FrameType_Error,
				RequestId:          unsubscribeFrame.RequestId,
				ConnectorChannelId: unsubscribeFrame.ConnectorChannelId,
				Code:               realtimetypes.ErrorCode_ChannelNotFound,
				Message:            "connectorChannelId is not subscribed on this connection",
			})
		}
		if err := g.leaseStore.ReleaseBlockPackSubscriber(
			channel.Id,
			fmt.Sprintf("%s:%d", connector.Id, unsubscribeFrame.ConnectorChannelId),
		); err != nil {
			logs.NotezyLogger.Error(ctx, err, "Failed to release realtime BlockPack subscriber lease")
		}

		g.workerManager.Detach(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_Detach,
			ChannelType:        channel.Type,
			ConnectionId:       connector.Id,
			ConnectorChannelId: unsubscribeFrame.ConnectorChannelId,
			ChannelId:          channel.Id,
		})
		metrics.NotezyMeter.Count(ctx, "realtime.channel.subscription.count", 1,
			attribute.String("action", "unsubscribe"),
			attribute.String("channelType", string(channel.Type)),
			attribute.String("outcome", "success"),
		)
		metrics.NotezyMeter.UpDown(ctx, "realtime.channel.count", -1,
			attribute.String("channelType", string(channel.Type)),
			attribute.String("permission", string(channel.Permission)),
		)

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
	switch frame.Type {
	case realtimetypes.InternalFrameType_LoadCompactableYjsDocument:
		input, err := g.yjsPersistenceService.GetCompactableYjsDocumentWithUpdates(
			context.Background(), frame.ChannelId,
		)
		if err != nil || input == nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsDocumentCompactionFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		payload, err := input.MarshalBytes()
		if err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsDocumentCompactionFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		g.workerManager.Forward(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_CompactableYjsDocumentLoaded,
			ChannelType:        frame.ChannelType,
			ConnectionId:       frame.ConnectionId,
			ConnectorChannelId: frame.ConnectorChannelId,
			ChannelId:          frame.ChannelId,
			Payload:            payload,
		})

		return
	case realtimetypes.InternalFrameType_ApplyCompactedYjsDocument:
		var result realtimetypes.YjsCompactionResult
		if err := result.UnmarshalBytes(frame.Payload); err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsDocumentCompactionFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		applied, err := g.yjsPersistenceService.ApplyCompactedYjsDocument(
			context.Background(), frame.ChannelId, result,
		)
		if err != nil || !applied {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsDocumentCompactionFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		g.workerManager.Forward(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_YjsDocumentCompacted,
			ChannelType:        frame.ChannelType,
			ConnectionId:       frame.ConnectionId,
			ConnectorChannelId: frame.ConnectorChannelId,
			ChannelId:          frame.ChannelId,
			Payload:            realtimetypes.MarshalYjsUpdateSequence(result.CutoffSequence),
		})

		return
	case realtimetypes.InternalFrameType_LoadYjsDocument:
		state, err := g.yjsPersistenceService.LoadDocument(context.Background(), frame.ChannelId)
		if err != nil {
			failureType := realtimetypes.YjsPersistenceFailureType_Retryable
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failureType = realtimetypes.YjsPersistenceFailureType_Terminal
			}

			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsPersistenceFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
				Payload:            []byte{byte(failureType)},
			})

			return
		}

		payload, err := state.MarshalBytes()
		if err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsPersistenceFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		g.workerManager.Forward(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_YjsDocumentLoaded,
			ChannelType:        frame.ChannelType,
			ConnectionId:       frame.ConnectionId,
			ConnectorChannelId: frame.ConnectorChannelId,
			ChannelId:          frame.ChannelId,
			Payload:            payload,
		})

		return
	case realtimetypes.InternalFrameType_AppendYjsUpdate:
		originConnectionId := frame.ConnectionId
		updateSequence, err := g.yjsPersistenceService.AppendUpdate(
			context.Background(),
			frame.ChannelId,
			uuid.New(),
			&originConnectionId,
			frame.Payload,
		)
		if err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsPersistenceFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		g.workerManager.Forward(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_YjsUpdatePersisted,
			ChannelType:        frame.ChannelType,
			ConnectionId:       frame.ConnectionId,
			ConnectorChannelId: frame.ConnectorChannelId,
			ChannelId:          frame.ChannelId,
			Payload:            realtimetypes.MarshalYjsUpdateSequence(updateSequence),
		})

		return
	case realtimetypes.InternalFrameType_AppendYjsUpdateBatch:
		var batch realtimetypes.YjsPersistenceBatch
		if err := batch.UnmarshalBytes(frame.Payload); err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsPersistenceFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
				Payload:            []byte{byte(realtimetypes.YjsPersistenceFailureType_Terminal)},
			})

			return
		}

		updateSequence, err := g.yjsPersistenceService.AppendUpdate(
			context.Background(),
			frame.ChannelId,
			batch.PersistenceBatchId,
			batch.OriginConnectionId,
			batch.Payload,
		)
		if err != nil {
			failureType := realtimetypes.YjsPersistenceFailureType_Retryable
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failureType = realtimetypes.YjsPersistenceFailureType_Terminal
			}

			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_YjsPersistenceFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
				Payload:            []byte{byte(failureType)},
			})

			return
		}

		g.workerManager.Forward(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_YjsUpdatePersisted,
			ChannelType:        frame.ChannelType,
			ConnectionId:       frame.ConnectionId,
			ConnectorChannelId: frame.ConnectorChannelId,
			ChannelId:          frame.ChannelId,
			Payload:            realtimetypes.MarshalYjsUpdateSequence(updateSequence),
		})

		return
	case realtimetypes.InternalFrameType_ApplyBlockProjection:
		var input dtos.ApplyBlockProjectionInput
		if err := json.Unmarshal(frame.Payload, &input); err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_BlockProjectionFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		result, err := g.blockProjectionService.Apply(context.Background(), frame.ChannelId, input)
		if err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_BlockProjectionFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		payload, err := json.Marshal(result)
		if err != nil {
			g.workerManager.Forward(realtimetypes.InternalFrame{
				Version:            byte(constants.RealtimeWorkerProtocolVersion),
				Type:               realtimetypes.InternalFrameType_BlockProjectionFailed,
				ChannelType:        frame.ChannelType,
				ConnectionId:       frame.ConnectionId,
				ConnectorChannelId: frame.ConnectorChannelId,
				ChannelId:          frame.ChannelId,
			})

			return
		}

		g.workerManager.Forward(realtimetypes.InternalFrame{
			Version:            byte(constants.RealtimeWorkerProtocolVersion),
			Type:               realtimetypes.InternalFrameType_BlockProjectionApplied,
			ChannelType:        frame.ChannelType,
			ConnectionId:       frame.ConnectionId,
			ConnectorChannelId: frame.ConnectorChannelId,
			ChannelId:          frame.ChannelId,
			Payload:            payload,
		})

		return
	}

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
		if err := g.leaseStore.ReleaseBlockPackSubscriber(
			channel.Id,
			fmt.Sprintf("%s:%d", connector.Id, frame.ConnectorChannelId),
		); err != nil {
			logs.NotezyLogger.Error(context.Background(), err, "Failed to release realtime BlockPack subscriber lease")
		}

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
		outcome := "resync_required"
		if frame.Type == realtimetypes.InternalFrameType_PermissionRevoked {
			outcome = "permission_revoked"
		}
		metrics.NotezyMeter.Count(context.Background(), "realtime.channel.subscription.count", 1,
			attribute.String("action", "detach"),
			attribute.String("channelType", string(channel.Type)),
			attribute.String("outcome", outcome),
		)
		metrics.NotezyMeter.UpDown(context.Background(), "realtime.channel.count", -1,
			attribute.String("channelType", string(channel.Type)),
			attribute.String("permission", string(channel.Permission)),
		)

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

	if err := connector.writeBinary(realtimetypes.BinaryFrame{
		Version:            byte(constants.RealtimeProtocolVersion),
		Type:               binaryFrameType,
		ConnectorChannelId: frame.ConnectorChannelId,
		Payload:            frame.Payload,
	}); err != nil {
		g.handleChannelBackpressure(connector, channel)

		return
	}
	metrics.NotezyMeter.Count(context.Background(), "realtime.frame.count", 1,
		attribute.String("direction", "outbound"),
		attribute.String("channelType", string(channel.Type)),
		attribute.String("frameType", string(binaryFrameType)),
	)
	metrics.NotezyMeter.Bytes(context.Background(), "realtime.payload.bytes", int64(len(frame.Payload)),
		attribute.String("direction", "outbound"),
		attribute.String("channelType", string(channel.Type)),
		attribute.String("frameType", string(binaryFrameType)),
	)
}

// Backpressure means this connector cannot write this channel's queued frames to its client fast enough.
// Do not silently discard Yjs document updates: losing one would leave the client with an incomplete
// CRDT history. Instead, detach only the congested channel, clear its pending outbound queue, and stop
// worker fanout for that channel. The control error is then sent with priority so the client can
// resubscribe and receive a complete state without disrupting other channels on the same connection.
func (g *Gateway) handleChannelBackpressure(
	connector *Connector,
	channel realtimetypes.Channel,
) {
	metrics.NotezyMeter.Count(
		context.Background(),
		"realtime.channel.backpressure.count",
		1,
		attribute.String("channelType", string(channel.Type)),
	)

	connectorChannelId, exists := connector.findChannel(channel.Type, channel.Id)
	if !exists {
		return
	}

	connector.unsubscribe(connectorChannelId)
	if err := g.leaseStore.ReleaseBlockPackSubscriber(
		channel.Id,
		fmt.Sprintf("%s:%d", connector.Id, connectorChannelId),
	); err != nil {
		logs.NotezyLogger.Error(context.Background(), err, "Failed to release realtime BlockPack subscriber lease")
	}

	g.workerManager.Detach(realtimetypes.InternalFrame{
		Version:            byte(constants.RealtimeWorkerProtocolVersion),
		Type:               realtimetypes.InternalFrameType_Detach,
		ChannelType:        channel.Type,
		ConnectionId:       connector.Id,
		ConnectorChannelId: connectorChannelId,
		ChannelId:          channel.Id,
	})
	metrics.NotezyMeter.Count(context.Background(), "realtime.channel.subscription.count", 1,
		attribute.String("action", "detach"),
		attribute.String("channelType", string(channel.Type)),
		attribute.String("outcome", "backpressure"),
	)
	metrics.NotezyMeter.UpDown(context.Background(), "realtime.channel.count", -1,
		attribute.String("channelType", string(channel.Type)),
		attribute.String("permission", string(channel.Permission)),
	)

	if !connector.writeError(realtimetypes.ErrorFrame{
		Version:            constants.RealtimeProtocolVersion,
		Type:               realtimetypes.FrameType_Error,
		ChannelType:        channel.Type,
		ChannelId:          &channel.Id,
		ConnectorChannelId: connectorChannelId,
		Code:               realtimetypes.ErrorCode_ChannelBackpressure,
		Message:            "channel outbound queue is full; resubscribe this channel to resync",
	}) {
		_ = connector.connection.Close()
	}
}
