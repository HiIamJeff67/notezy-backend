package realtimetypes

import sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"

type ErrorCode = sharedtypes.ErrorCode

const (
	ErrorCode_AuthenticationManagedByUpgrade = sharedtypes.ErrorCode_AuthenticationManagedByUpgrade
	ErrorCode_BinaryChannelNotReady          = sharedtypes.ErrorCode_BinaryChannelNotReady
	ErrorCode_ChannelLimitExceeded           = sharedtypes.ErrorCode_ChannelLimitExceeded
	ErrorCode_ChannelNotFound                = sharedtypes.ErrorCode_ChannelNotFound
	ErrorCode_ChannelPermissionDenied        = sharedtypes.ErrorCode_ChannelPermissionDenied
	ErrorCode_ChannelBackpressure            = sharedtypes.ErrorCode_ChannelBackpressure
	ErrorCode_InvalidAcknowledgement         = sharedtypes.ErrorCode_InvalidAcknowledgement
	ErrorCode_InvalidBinaryFrame             = sharedtypes.ErrorCode_InvalidBinaryFrame
	ErrorCode_InvalidChannelId               = sharedtypes.ErrorCode_InvalidChannelId
	ErrorCode_InvalidChannelTicket           = sharedtypes.ErrorCode_InvalidChannelTicket
	ErrorCode_InvalidChannelType             = sharedtypes.ErrorCode_InvalidChannelType
	ErrorCode_InvalidConnectorChannelId      = sharedtypes.ErrorCode_InvalidConnectorChannelId
	ErrorCode_InvalidControlFrame            = sharedtypes.ErrorCode_InvalidControlFrame
	ErrorCode_PermissionRevoked              = sharedtypes.ErrorCode_PermissionRevoked
	ErrorCode_ResourceUnavailable            = sharedtypes.ErrorCode_ResourceUnavailable
	ErrorCode_RoomAdmissionUnavailable       = sharedtypes.ErrorCode_RoomAdmissionUnavailable
	ErrorCode_RoomConnectionLimitExceeded    = sharedtypes.ErrorCode_RoomConnectionLimitExceeded
	ErrorCode_ResubscribeRequired            = sharedtypes.ErrorCode_ResubscribeRequired
	ErrorCode_UnsupportedBinaryType          = sharedtypes.ErrorCode_UnsupportedBinaryType
	ErrorCode_UnsupportedChannelType         = sharedtypes.ErrorCode_UnsupportedChannelType
	ErrorCode_UnsupportedControlType         = sharedtypes.ErrorCode_UnsupportedControlType
	ErrorCode_UnsupportedMessageType         = sharedtypes.ErrorCode_UnsupportedMessageType
	ErrorCode_UnsupportedProtocolVersion     = sharedtypes.ErrorCode_UnsupportedProtocolVersion
	ErrorCode_WorkerUnavailable              = sharedtypes.ErrorCode_WorkerUnavailable
)
