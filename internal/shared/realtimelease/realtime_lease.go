package realtimelease

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

type RealtimeBlockPackChannelRevocation struct {
	UserId       uuid.UUID   `json:"userId"`
	BlockPackIds []uuid.UUID `json:"blockPackIds"`
}

type RealtimeLeaseStore struct {
	Range types.Range[int, int]
}

type RealtimeLeaseCacheStore struct {
	databaseNumber int
	redisClient    *redis.Client
}

func NewRealtimeLeaseStore() *RealtimeLeaseStore {
	return &RealtimeLeaseStore{
		Range: types.Range[int, int]{
			Start: constants.RealtimeRedisServerNumber,
			Size:  1,
		},
	}
}

func NewRealtimeLeaseCacheStore(
	databaseNumber int,
	redisClient *redis.Client,
) *RealtimeLeaseCacheStore {
	return &RealtimeLeaseCacheStore{
		databaseNumber: databaseNumber,
		redisClient:    redisClient,
	}
}

func Register(
	ctx context.Context,
	clientManager *platformredis.ClientManager,
) error {
	leaseStore := NewRealtimeLeaseStore()
	if err := clientManager.ConnectAll(leaseStore.Range); err != nil {
		return err
	}

	redisClient, exists := clientManager.Client(constants.RealtimeRedisServerNumber)
	if !exists {
		return errors.New("realtime Redis client is unavailable")
	}

	return platformredis.RegisterCacheStores(
		ctx,
		NewRealtimeLeaseCacheStore(constants.RealtimeRedisServerNumber, redisClient),
	)
}

func (s *RealtimeLeaseCacheStore) DatabaseNumber() int {
	return s.databaseNumber
}

func (s *RealtimeLeaseCacheStore) Initialize(_ context.Context) error {
	return nil
}

func (s *RealtimeLeaseStore) getRedisClient() (*redis.Client, error) {
	cacheStore, err := platformredis.GetRedisCacheStore(constants.RealtimeRedisServerNumber)
	if err != nil {
		return nil, errors.New("realtime Redis lease store is unavailable")
	}
	realtimeLeaseCacheStore, ok := cacheStore.(*RealtimeLeaseCacheStore)
	if !ok || realtimeLeaseCacheStore.redisClient == nil {
		return nil, errors.New("realtime Redis lease store is unavailable")
	}

	return realtimeLeaseCacheStore.redisClient, nil
}

func (s *RealtimeLeaseStore) blockPackChannelRevocationKey() string {
	return "Realtime:blockPack:channel-revocations"
}

func (s *RealtimeLeaseStore) PublishBlockPackChannelRevocation(
	userId uuid.UUID,
	blockPackIds []uuid.UUID,
) error {
	if len(blockPackIds) == 0 {
		return nil
	}

	redisClient, err := s.getRedisClient()
	if err != nil {
		return err
	}

	payload, err := json.Marshal(RealtimeBlockPackChannelRevocation{
		UserId:       userId,
		BlockPackIds: blockPackIds,
	})
	if err != nil {
		return err
	}

	return redisClient.Publish(s.blockPackChannelRevocationKey(), payload).Err()
}
