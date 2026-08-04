package workers

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/internal/realtimegateway/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type realtimeWorker struct {
	endpoint string

	activeChannels      map[string]realtimetypes.InternalFrame
	activeChannelsMutex sync.RWMutex

	frameHandler      func(realtimetypes.InternalFrame)
	frameHandlerMutex sync.RWMutex

	outbound        chan realtimetypes.InternalFrame
	ready           atomic.Bool
	connectionMutex sync.Mutex
	connection      *websocket.Conn
}

func (w *realtimeWorker) channelKey(frame realtimetypes.InternalFrame) string {
	return frame.ConnectionId.String() + ":" + strconv.FormatUint(uint64(frame.ConnectorChannelId), 10)
}

func (w *realtimeWorker) enqueue(frame realtimetypes.InternalFrame) bool {
	select {
	case w.outbound <- frame:
		return true
	default:
		return false
	}
}

func (w *realtimeWorker) run(ctx context.Context) {
	for {
		connection, _, err := websocket.DefaultDialer.DialContext(ctx, w.endpoint, nil)
		if err != nil {
			w.ready.Store(false)
			if !wait(ctx, constants.RealtimeWorkerReconnectDelay) {
				return
			}

			continue
		}
		w.connectionMutex.Lock()
		w.connection = connection
		w.connectionMutex.Unlock()

		w.ready.Store(true)
		if !w.replayActiveChannels(connection) {
			w.ready.Store(false)
			connection.Close()
			w.connectionMutex.Lock()
			w.connection = nil
			w.connectionMutex.Unlock()
			if !wait(ctx, constants.RealtimeWorkerReconnectDelay) {
				return
			}

			continue
		}

		readError := make(chan struct{})
		go w.read(connection, readError)

		connected := true
		for connected {
			select {
			case <-ctx.Done():
				connected = false
			case <-readError:
				connected = false
			case frame := <-w.outbound:
				payload, err := frame.MarshalBytes()
				if err != nil || connection.SetWriteDeadline(time.Now().Add(constants.RealtimeControlWriteTimeout)) != nil ||
					connection.WriteMessage(websocket.BinaryMessage, payload) != nil {
					connected = false
				}
			}
		}

		w.ready.Store(false)
		connection.Close()
		w.connectionMutex.Lock()
		w.connection = nil
		w.connectionMutex.Unlock()
		if !wait(ctx, constants.RealtimeWorkerReconnectDelay) {
			return
		}
	}
}

func (w *realtimeWorker) close() {
	w.connectionMutex.Lock()
	defer w.connectionMutex.Unlock()

	if w.connection != nil {
		_ = w.connection.Close()
	}
}

func wait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (w *realtimeWorker) replayActiveChannels(connection *websocket.Conn) bool {
	w.activeChannelsMutex.RLock()
	frames := make([]realtimetypes.InternalFrame, 0, len(w.activeChannels))
	for _, frame := range w.activeChannels {
		frames = append(frames, frame)
	}
	w.activeChannelsMutex.RUnlock()

	for _, frame := range frames {
		payload, err := frame.MarshalBytes()
		if err != nil || connection.SetWriteDeadline(time.Now().Add(constants.RealtimeControlWriteTimeout)) != nil ||
			connection.WriteMessage(websocket.BinaryMessage, payload) != nil {
			return false
		}
	}

	return true
}

func (w *realtimeWorker) read(connection *websocket.Conn, readError chan<- struct{}) {
	defer close(readError)

	for {
		messageType, payload, err := connection.ReadMessage()
		if err != nil || messageType != websocket.BinaryMessage {
			return
		}

		var frame realtimetypes.InternalFrame
		if err := frame.UnmarshalBytes(payload); err != nil || int(frame.Version) != constants.RealtimeWorkerProtocolVersion {
			return
		}

		w.frameHandlerMutex.RLock()
		frameHandler := w.frameHandler
		w.frameHandlerMutex.RUnlock()

		if frameHandler != nil {
			frameHandler(frame)
		}
	}
}
