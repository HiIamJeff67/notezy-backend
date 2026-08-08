package redis

import (
	"errors"
	"fmt"
	"hash/fnv"

	redisclient "github.com/go-redis/redis"
)

// ClientSet owns the Redis clients used by one runtime. It is immutable after
// construction, so concurrent lookups do not require a lock.
type ClientSet struct {
	clients []*redisclient.Client
}

func NewClientSet(config Config) (*ClientSet, error) {
	client := redisclient.NewClient(&redisclient.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.Database,
	})
	if _, err := client.Ping().Result(); err != nil {
		_ = client.Close()

		return nil, fmt.Errorf("connect Redis database %d: %w", config.Database, err)
	}

	return NewClientSetFromClients(client), nil
}

func NewClientSetFromClients(clients ...*redisclient.Client) *ClientSet {
	ownedClients := make([]*redisclient.Client, len(clients))
	copy(ownedClients, clients)

	return &ClientSet{clients: ownedClients}
}

func (s *ClientSet) Clients() []*redisclient.Client {
	if s == nil {
		return nil
	}

	clients := make([]*redisclient.Client, len(s.clients))
	copy(clients, s.clients)

	return clients
}

func (s *ClientSet) Client(index int) (*redisclient.Client, bool) {
	if s == nil || index < 0 || index >= len(s.clients) {
		return nil, false
	}

	return s.clients[index], true
}

func (s *ClientSet) ClientForKey(key string) (*redisclient.Client, int, error) {
	if s == nil || len(s.clients) == 0 {
		return nil, 0, errors.New("Redis client set is unavailable")
	}

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	index := int(hash.Sum32()) % len(s.clients)

	return s.clients[index], index, nil
}

func (s *ClientSet) Close() error {
	if s == nil {
		return nil
	}

	var exceptions []error
	for _, client := range s.clients {
		if client == nil {
			continue
		}
		if err := client.Close(); err != nil {
			exceptions = append(exceptions, err)
		}
	}

	return errors.Join(exceptions...)
}
