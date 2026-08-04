package redis

import "sync"

var defaultClientManagerMutex sync.RWMutex
var DefaultClientManager *ClientManager

func InitializeDefaultClientManager(config Config) *ClientManager {
	defaultClientManagerMutex.Lock()
	defer defaultClientManagerMutex.Unlock()

	DefaultClientManager = NewClientManager(config)

	return DefaultClientManager
}
