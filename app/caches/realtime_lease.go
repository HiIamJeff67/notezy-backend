package caches

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	redisscripts "github.com/HiIamJeff67/notezy-backend/app/caches/scripts"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
)

type RealtimeLeaseStore struct {
	redisClient *redis.Client
}

func NewRealtimeLeaseStore() *RealtimeLeaseStore {
	return &RealtimeLeaseStore{
		redisClient: RedisClientMap[RealtimeRange.Start],
	}
}

func NewRealtimeLeaseStoreWithClient(redisClient *redis.Client) *RealtimeLeaseStore {
	return &RealtimeLeaseStore{
		redisClient: redisClient,
	}
}

func (s *RealtimeLeaseStore) AcquireUserConnection(
	userPublicId uuid.UUID,
	connectionId uuid.UUID,
	maximumConnections int,
) (bool, int64, error) {
	return s.acquire(
		fmt.Sprintf("Realtime:user:%s:connections", userPublicId),
		connectionId.String(),
		maximumConnections,
	)
}

func (s *RealtimeLeaseStore) RefreshUserConnection(
	userPublicId uuid.UUID,
	connectionId uuid.UUID,
) (bool, error) {
	return s.refresh(
		fmt.Sprintf("Realtime:user:%s:connections", userPublicId),
		connectionId.String(),
	)
}

func (s *RealtimeLeaseStore) ReleaseUserConnection(
	userPublicId uuid.UUID,
	connectionId uuid.UUID,
) error {
	return s.release(
		fmt.Sprintf("Realtime:user:%s:connections", userPublicId),
		connectionId.String(),
	)
}

func (s *RealtimeLeaseStore) AcquireBlockPackSubscriber(
	blockPackId uuid.UUID,
	member string,
	maximumSubscribers int,
) (bool, int64, error) {
	return s.acquire(
		fmt.Sprintf("Realtime:blockPack:%s:subscribers", blockPackId),
		member,
		maximumSubscribers,
	)
}

func (s *RealtimeLeaseStore) RefreshBlockPackSubscriber(
	blockPackId uuid.UUID,
	member string,
) (bool, error) {
	return s.refresh(
		fmt.Sprintf("Realtime:blockPack:%s:subscribers", blockPackId),
		member,
	)
}

func (s *RealtimeLeaseStore) ReleaseBlockPackSubscriber(
	blockPackId uuid.UUID,
	member string,
) error {
	return s.release(
		fmt.Sprintf("Realtime:blockPack:%s:subscribers", blockPackId), member)
}

func (s *RealtimeLeaseStore) acquire(key string, member string, maximumMembers int) (bool, int64, error) {
	if s == nil || s.redisClient == nil {
		return false, 0, errors.New("realtime redis lease store is unavailable")
	}
	if maximumMembers <= 0 {
		return false, 0, nil
	}

	now := time.Now()
	result, err := redisscripts.AcquireRealtimeLease.Eval(
		s.redisClient,
		[]string{key},
		now.UnixMilli(),
		now.Add(constants.RealtimeLeaseTTL).UnixMilli(),
		maximumMembers,
		constants.RealtimeLeaseTTL.Milliseconds(),
		member,
	).Result()
	if err != nil {
		return false, 0, err
	}

	values, ok := result.([]interface{})
	if !ok || len(values) != 2 {
		return false, 0, errors.New("realtime redis lease acquisition returned an invalid result")
	}

	acquired, ok := values[0].(int64)
	if !ok {
		return false, 0, errors.New("realtime redis lease acquisition returned an invalid status")
	}
	activeMembers, ok := values[1].(int64)
	if !ok {
		return false, 0, errors.New("realtime redis lease acquisition returned an invalid member count")
	}

	return acquired == 1, activeMembers, nil
}

func (s *RealtimeLeaseStore) refresh(key string, member string) (bool, error) {
	if s == nil || s.redisClient == nil {
		return false, errors.New("realtime redis lease store is unavailable")
	}

	now := time.Now()
	result, err := redisscripts.RefreshRealtimeLease.Eval(
		s.redisClient,
		[]string{key},
		now.UnixMilli(),
		now.Add(constants.RealtimeLeaseTTL).UnixMilli(),
		constants.RealtimeLeaseTTL.Milliseconds(),
		member,
	).Result()
	if err != nil {
		return false, err
	}

	refreshed, ok := result.(int64)
	if !ok {
		return false, errors.New("realtime redis lease refresh returned an invalid result")
	}

	return refreshed == 1, nil
}

func (s *RealtimeLeaseStore) release(key string, member string) error {
	if s == nil || s.redisClient == nil {
		return errors.New("realtime redis lease store is unavailable")
	}

	return redisscripts.ReleaseRealtimeLease.Eval(s.redisClient, []string{key}, member).Err()
}
