package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"

	redisclient "github.com/go-redis/redis"

	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type ClientManager struct {
	config Config

	clients        map[int]*redisclient.Client
	clientToConfig map[*redisclient.Client]Config
	clientMapMutex sync.Mutex
}

/* ============================== Constructor ============================== */

func NewClientManager(config Config) *ClientManager {
	return &ClientManager{
		config: config,

		clients:        make(map[int]*redisclient.Client),
		clientToConfig: make(map[*redisclient.Client]Config),
	}
}

/* ============================== Getter Methods ============================== */

func (m *ClientManager) Clients() map[int]*redisclient.Client {
	return m.clients
}

func (m *ClientManager) Client(databaseNumber int) (*redisclient.Client, bool) {
	m.clientMapMutex.Lock()
	defer m.clientMapMutex.Unlock()

	client, exists := m.clients[databaseNumber]

	return client, exists
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
	config.Database = serverNumber
	client := redisclient.NewClient(&redisclient.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.Database,
	})
	if _, err := client.Ping().Result(); err != nil {
		_ = client.Close()

		return nil, fmt.Errorf("connect Redis database %d: %w", config.Database, err)
	}

	m.clients[config.Database] = client
	m.clientToConfig[client] = config
	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Redis client server of %d connected", config.Database))

	return client, nil
}

func (m *ClientManager) DisconnectAll() error {
	m.clientMapMutex.Lock()
	clients := make(map[*redisclient.Client]Config, len(m.clientToConfig))
	for client, config := range m.clientToConfig {
		clients[client] = config
	}
	m.clientMapMutex.Unlock()

	var exceptions []error
	for client, config := range clients {
		if err := client.Close(); err != nil {
			exceptions = append(exceptions, fmt.Errorf("disconnect Redis database %d: %w", config.Database, err))
			continue
		}

		m.clientMapMutex.Lock()
		delete(m.clientToConfig, client)
		delete(m.clients, config.Database)
		m.clientMapMutex.Unlock()
		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Redis client server of %d disconnected", config.Database))
	}

	return errors.Join(exceptions...)
}
