package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"

	redisclient "github.com/go-redis/redis"

	configs "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type ClientManager struct {
	config configs.CacheManagerConfig

	clients        map[int]*redisclient.Client
	clientToConfig map[*redisclient.Client]configs.CacheManagerConfig
	clientMapMutex sync.Mutex
}

/* ============================== Constructor ============================== */

func NewClientManager(config configs.CacheManagerConfig) *ClientManager {
	return &ClientManager{
		config: config,

		clients:        make(map[int]*redisclient.Client),
		clientToConfig: make(map[*redisclient.Client]configs.CacheManagerConfig),
	}
}

/* ============================== Getter Methods ============================== */

func (m *ClientManager) Clients() map[int]*redisclient.Client {
	return m.clients
}

/* ============================== Connection Methods ============================== */

func (m *ClientManager) ConnectAll(serverRanges ...types.Range[int, int]) error {
	for _, serverRange := range serverRanges {
		for serverNumber := serverRange.Start; serverNumber < serverRange.Start+serverRange.Size; serverNumber++ {
			if _, err := m.Connect(serverNumber); err != nil {
				_ = m.DisconnectAll()

				return err
			}
		}
	}

	return nil
}

func (m *ClientManager) Connect(serverNumber int) (*redisclient.Client, error) {
	m.clientMapMutex.Lock()
	defer m.clientMapMutex.Unlock()

	if client, exists := m.clients[serverNumber]; exists {
		return client, nil
	}

	config := m.config
	config.DB = serverNumber
	client := redisclient.NewClient(&redisclient.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})
	if _, err := client.Ping().Result(); err != nil {
		_ = client.Close()

		return nil, fmt.Errorf("connect Redis database %d: %w", config.DB, err)
	}

	m.clients[config.DB] = client
	m.clientToConfig[client] = config
	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Redis client server of %d connected", config.DB))

	return client, nil
}

func (m *ClientManager) DisconnectAll() error {
	m.clientMapMutex.Lock()
	clients := make(map[*redisclient.Client]configs.CacheManagerConfig, len(m.clientToConfig))
	for client, config := range m.clientToConfig {
		clients[client] = config
	}
	m.clientMapMutex.Unlock()

	var exceptions []error
	for client, config := range clients {
		if err := client.Close(); err != nil {
			exceptions = append(exceptions, fmt.Errorf("disconnect Redis database %d: %w", config.DB, err))
			continue
		}

		m.clientMapMutex.Lock()
		delete(m.clientToConfig, client)
		delete(m.clients, config.DB)
		m.clientMapMutex.Unlock()
		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Redis client server of %d disconnected", config.DB))
	}

	return errors.Join(exceptions...)
}
