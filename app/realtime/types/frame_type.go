package realtimetypes

type FrameType string

const (
	FrameType_Ready        FrameType = "ready"
	FrameType_Error        FrameType = "error"
	FrameType_Authenticate FrameType = "authenticate"
	FrameType_Ping         FrameType = "ping"
	FrameType_Pong         FrameType = "pong"
	FrameType_Subscribe    FrameType = "subscribe"
	FrameType_Subscribed   FrameType = "subscribed"
	FrameType_Unsubscribe  FrameType = "unsubscribe"
	FrameType_Unsubscribed FrameType = "unsubscribed"
	FrameType_Acknowledge  FrameType = "ack"
	FrameType_Acknowledged FrameType = "acknowledged"
	FrameType_Heartbeat    FrameType = "heartbeat"
	FrameType_Reconnect    FrameType = "reconnect"
)

type ErrorCode string

const (
	ErrorCode_AuthenticationManagedByUpgrade ErrorCode = "authentication_managed_by_upgrade"
	ErrorCode_BinaryChannelNotReady          ErrorCode = "binary_channel_not_ready"
	ErrorCode_ChannelLimitExceeded           ErrorCode = "channel_limit_exceeded"
	ErrorCode_ChannelNotFound                ErrorCode = "channel_not_found"
	ErrorCode_ChannelPermissionDenied        ErrorCode = "channel_permission_denied"
	ErrorCode_InvalidAcknowledgement         ErrorCode = "invalid_acknowledgement"
	ErrorCode_InvalidBinaryFrame             ErrorCode = "invalid_binary_frame"
	ErrorCode_InvalidChannelId               ErrorCode = "invalid_channel_id"
	ErrorCode_InvalidChannelTicket           ErrorCode = "invalid_channel_ticket"
	ErrorCode_InvalidChannelType             ErrorCode = "invalid_channel_type"
	ErrorCode_InvalidConnectorChannelId      ErrorCode = "invalid_connector_channel_id"
	ErrorCode_InvalidControlFrame            ErrorCode = "invalid_control_frame"
	ErrorCode_PermissionRevoked              ErrorCode = "permission_revoked"
	ErrorCode_ResubscribeRequired            ErrorCode = "resubscribe_required"
	ErrorCode_UnsupportedBinaryType          ErrorCode = "unsupported_binary_type"
	ErrorCode_UnsupportedChannelType         ErrorCode = "unsupported_channel_type"
	ErrorCode_UnsupportedControlType         ErrorCode = "unsupported_control_type"
	ErrorCode_UnsupportedMessageType         ErrorCode = "unsupported_message_type"
	ErrorCode_UnsupportedProtocolVersion     ErrorCode = "unsupported_protocol_version"
	ErrorCode_WorkerUnavailable              ErrorCode = "worker_unavailable"
)

type BinaryFrameType byte

const (
	BinaryFrameType_YjsDocument BinaryFrameType = 1
	BinaryFrameType_Awareness   BinaryFrameType = 2
)

type InternalFrameType byte

const (
	InternalFrameType_Attach                 InternalFrameType = 1
	InternalFrameType_Detach                 InternalFrameType = 2
	InternalFrameType_YjsDocument            InternalFrameType = 3
	InternalFrameType_Awareness              InternalFrameType = 4
	InternalFrameType_ResyncRequired         InternalFrameType = 5
	InternalFrameType_PermissionRevoked      InternalFrameType = 6
	InternalFrameType_LoadYjsDocument        InternalFrameType = 7
	InternalFrameType_YjsDocumentLoaded      InternalFrameType = 8
	InternalFrameType_AppendYjsUpdate        InternalFrameType = 9
	InternalFrameType_YjsUpdatePersisted     InternalFrameType = 10
	InternalFrameType_YjsPersistenceFailed   InternalFrameType = 11
	InternalFrameType_ApplyBlockProjection   InternalFrameType = 12
	InternalFrameType_BlockProjectionApplied InternalFrameType = 13
	InternalFrameType_BlockProjectionFailed  InternalFrameType = 14
)

type InternalChannelType byte

const (
	InternalChannelType_BlockPack InternalChannelType = 1
)
