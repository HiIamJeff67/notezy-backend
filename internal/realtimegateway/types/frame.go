package realtimetypes

import "github.com/google/uuid"

type ControlFrame struct {
	Version   int       `json:"version"`
	Type      FrameType `json:"type"`
	RequestId string    `json:"requestId,omitempty"`
}

type ReadyFrame struct {
	Version             int       `json:"version"`
	Type                FrameType `json:"type"`
	ConnectionId        string    `json:"connectionId"`
	ResubscribeRequired bool      `json:"resubscribeRequired"`
}

type ErrorFrame struct {
	Version            int         `json:"version"`
	Type               FrameType   `json:"type"`
	RequestId          string      `json:"requestId,omitempty"`
	ChannelType        ChannelType `json:"channelType,omitempty"`
	ChannelId          *uuid.UUID  `json:"channelId,omitempty"`
	ConnectorChannelId uint32      `json:"connectorChannelId,omitempty"`
	Code               ErrorCode   `json:"code"`
	Message            string      `json:"message"`
}

type SubscribeFrame struct {
	Version       int         `json:"version"`
	Type          FrameType   `json:"type"`
	RequestId     string      `json:"requestId,omitempty"`
	ChannelType   ChannelType `json:"channelType"`
	ChannelId     uuid.UUID   `json:"channelId"`
	ChannelTicket string      `json:"channelTicket,omitempty"`
}

type SubscribedFrame struct {
	Version            int                            `json:"version"`
	Type               FrameType                      `json:"type"`
	RequestId          string                         `json:"requestId,omitempty"`
	ChannelType        ChannelType                    `json:"channelType"`
	ChannelId          uuid.UUID                      `json:"channelId"`
	ConnectorChannelId uint32                         `json:"connectorChannelId"`
	Existing           bool                           `json:"existing"`
	Participants       []BlockPackPresenceParticipant `json:"participants,omitempty"`
}

type BlockPackPresenceParticipant struct {
	UserPublicId      uuid.UUID `json:"userPublicId"`
	ChannelPermission string    `json:"channelPermission"`
	ConnectionCount   int       `json:"connectionCount"`
}

type BlockPackPresenceDeltaFrame struct {
	Version     int                          `json:"version"`
	Type        FrameType                    `json:"type"`
	ChannelType ChannelType                  `json:"channelType"`
	ChannelId   uuid.UUID                    `json:"channelId"`
	Participant BlockPackPresenceParticipant `json:"participant"`
}

type UnsubscribeFrame struct {
	Version            int       `json:"version"`
	Type               FrameType `json:"type"`
	RequestId          string    `json:"requestId,omitempty"`
	ConnectorChannelId uint32    `json:"connectorChannelId"`
}

type UnsubscribedFrame struct {
	Version            int         `json:"version"`
	Type               FrameType   `json:"type"`
	RequestId          string      `json:"requestId,omitempty"`
	ChannelType        ChannelType `json:"channelType"`
	ChannelId          uuid.UUID   `json:"channelId"`
	ConnectorChannelId uint32      `json:"connectorChannelId"`
}

type AcknowledgeFrame struct {
	Version            int       `json:"version"`
	Type               FrameType `json:"type"`
	RequestId          string    `json:"requestId,omitempty"`
	ConnectorChannelId uint32    `json:"connectorChannelId"`
	Sequence           int64     `json:"sequence"`
}

type AcknowledgedFrame struct {
	Version            int       `json:"version"`
	Type               FrameType `json:"type"`
	RequestId          string    `json:"requestId,omitempty"`
	ConnectorChannelId uint32    `json:"connectorChannelId"`
	Sequence           int64     `json:"sequence"`
}

type HeartbeatFrame struct {
	Version      int       `json:"version"`
	Type         FrameType `json:"type"`
	RequestId    string    `json:"requestId,omitempty"`
	UnixMilliNow int64     `json:"unixMilliNow"`
}

type ResourceEventFrame struct {
	Version            int        `json:"version"`
	Type               FrameType  `json:"type"`
	EventId            uuid.UUID  `json:"eventId"`
	EventType          string     `json:"eventType"`
	ResourceId         uuid.UUID  `json:"resourceId"`
	TargetUserPublicId *uuid.UUID `json:"targetUserPublicId,omitempty"`
	Change             string     `json:"change"`
	Permission         string     `json:"permission,omitempty"`
}
