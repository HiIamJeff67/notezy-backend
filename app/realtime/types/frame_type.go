package realtimetypes

// The basic frame type for websocket connection
// note that there are two types of binary frame and internal frame
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
