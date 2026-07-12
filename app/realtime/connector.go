package realtime

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Connector struct {
	Id            uuid.UUID
	UserPublicId  uuid.UUID
	UserAgent     string
	connection    *websocket.Conn
	channels      map[uint32]realtimetypes.Channel
	nextChannelId uint32
	channelMutex  sync.RWMutex
	writeMutex    sync.Mutex
}

func (c *Connector) get(connectorChannelId uint32) (realtimetypes.Channel, bool) {
	c.channelMutex.RLock()
	defer c.channelMutex.RUnlock()

	channel, exists := c.channels[connectorChannelId]

	return channel, exists
}

func (c *Connector) subscribe(channel realtimetypes.Channel) (uint32, bool) {
	c.channelMutex.Lock()
	defer c.channelMutex.Unlock()

	for connectorChannelId, existingChannel := range c.channels {
		if existingChannel.Type == channel.Type && existingChannel.Id == channel.Id {
			return connectorChannelId, true
		}
	}
	if len(c.channels) >= constants.RealtimeMaxChannelsPerConnection || c.nextChannelId == constants.MAX_UINT32 {
		return 0, false
	}

	c.nextChannelId++
	c.channels[c.nextChannelId] = channel

	return c.nextChannelId, false
}

func (c *Connector) unsubscribe(connectorChannelId uint32) (realtimetypes.Channel, bool) {
	c.channelMutex.Lock()
	defer c.channelMutex.Unlock()

	channel, exists := c.channels[connectorChannelId]
	if exists {
		delete(c.channels, connectorChannelId)
	}

	return channel, exists
}

func (c *Connector) acknowledge(connectorChannelId uint32, sequence int64) (bool, bool) {
	c.channelMutex.Lock()
	defer c.channelMutex.Unlock()

	channel, exists := c.channels[connectorChannelId]
	if !exists {
		return false, false
	}
	if sequence < channel.AcknowledgedSequence {
		return true, false
	}

	channel.AcknowledgedSequence = sequence
	c.channels[connectorChannelId] = channel

	return true, true
}

func (c *Connector) writeError(frame realtimetypes.ErrorFrame) bool {
	return c.writeJSON(frame) == nil
}

func (c *Connector) writeJSON(frame any) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	if err := c.connection.SetWriteDeadline(time.Now().Add(constants.RealtimeControlWriteTimeout)); err != nil {
		return err
	}

	return c.connection.WriteJSON(frame)
}

func (c *Connector) writeControl(messageType int, payload []byte) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	return c.connection.WriteControl(
		messageType,
		payload,
		time.Now().Add(constants.RealtimeControlWriteTimeout),
	)
}

func (c *Connector) writeBinary(payload []byte) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()

	if err := c.connection.SetWriteDeadline(time.Now().Add(constants.RealtimeControlWriteTimeout)); err != nil {
		return err
	}

	return c.connection.WriteMessage(websocket.BinaryMessage, payload)
}
