package realtime

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	realtimetypes "github.com/HiIamJeff67/notezy-backend/app/realtime/types"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type Channel struct {
	Type                 realtimetypes.ChannelType
	Id                   uuid.UUID
	acknowledgedSequence int64
}

type Connector struct {
	connection    *websocket.Conn
	channels      map[uint32]Channel
	nextChannelId uint32
	channelMutex  sync.RWMutex
	writeMutex    sync.Mutex
}

func (c *Connector) get(connectorChannelId uint32) (Channel, bool) {
	c.channelMutex.RLock()
	defer c.channelMutex.RUnlock()

	channel, exists := c.channels[connectorChannelId]

	return channel, exists
}

func (c *Connector) append(channel Channel) (uint32, bool) {
	c.channelMutex.Lock()
	defer c.channelMutex.Unlock()

	for connectorChannelId, existingChannel := range c.channels {
		if existingChannel.Type == channel.Type && existingChannel.Id == channel.Id {
			return connectorChannelId, true
		}
	}
	if len(c.channels) >= constants.RealtimeMaxChannelsPerConnection || c.nextChannelId == ^uint32(0) {
		return 0, false
	}

	c.nextChannelId++
	c.channels[c.nextChannelId] = channel

	return c.nextChannelId, false
}

func (c *Connector) remove(connectorChannelId uint32) (Channel, bool) {
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
	if sequence < channel.acknowledgedSequence {
		return true, false
	}

	channel.acknowledgedSequence = sequence
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
