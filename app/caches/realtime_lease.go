package caches

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	redisscripts "github.com/HiIamJeff67/notezy-backend/app/caches/scripts"
	constants "github.com/HiIamJeff67/notezy-backend/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/shared/types"
)

type RealtimeLeaseStore struct {
	redisClientMap map[int]*redis.Client

	Range           types.Range[int, int]
	MaxServerNumber int
}

type RealtimeBlockPackParticipant struct {
	Member            string `json:"member"`
	UserPublicId      string `json:"userPublicId"`
	ChannelPermission string `json:"channelPermission"`
}

/* ============================== Constructor ============================== */

func NewRealtimeLeaseStore(redisClientMap map[int]*redis.Client) *RealtimeLeaseStore {
	rangeValue := types.Range[int, int]{Start: constants.RealtimeRedisServerNumber, Size: 1}

	return &RealtimeLeaseStore{
		redisClientMap:  redisClientMap,
		Range:           rangeValue,
		MaxServerNumber: rangeValue.Start + rangeValue.Size - 1,
	}
}

/* ============================== Auxiliary Methods ============================== */

func (s *RealtimeLeaseStore) getRedisClient(identifier string) (*redis.Client, error) {
	if s == nil || s.redisClientMap == nil || s.Range.Size <= 0 {
		return nil, errors.New("realtime redis lease store is unavailable")
	}

	serverNumber := min(s.MaxServerNumber, s.Range.Start+s.hashIdentifier(identifier))
	redisClient, ok := s.redisClientMap[serverNumber]
	if !ok || redisClient == nil {
		return nil, errors.New("realtime redis lease store is unavailable")
	}

	return redisClient, nil
}

func (s *RealtimeLeaseStore) hashIdentifier(identifier string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(identifier))

	return int(hash.Sum32()) % s.Range.Size
}

func (s *RealtimeLeaseStore) acquire(identifier string, key string, member string, maximumMembers int) (bool, int64, error) {
	redisClient, err := s.getRedisClient(identifier)
	if err != nil {
		return false, 0, err
	}

	if maximumMembers <= 0 {
		return false, 0, nil
	}

	now := time.Now()
	result, err := redisscripts.AcquireRealtimeLease.Eval(
		redisClient,
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

func (s *RealtimeLeaseStore) refresh(identifier string, key string, member string) (bool, error) {
	redisClient, err := s.getRedisClient(identifier)
	if err != nil {
		return false, err
	}

	now := time.Now()
	result, err := redisscripts.RefreshRealtimeLease.Eval(
		redisClient,
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

func (s *RealtimeLeaseStore) blockPackParticipantKey(blockPackId uuid.UUID) string {
	return fmt.Sprintf("Realtime:blockPack:%s:participants", blockPackId)
}

func (s *RealtimeLeaseStore) release(identifier string, key string, member string) error {
	redisClient, err := s.getRedisClient(identifier)
	if err != nil {
		return err
	}

	return redisscripts.ReleaseRealtimeLease.Eval(redisClient, []string{key}, member).Err()
}

/* ============================== User Connection Methods ============================== */

func (s *RealtimeLeaseStore) AcquireUserConnection(
	identifier uuid.UUID,
	connectionId uuid.UUID,
	maximumConnections int,
) (bool, int64, error) {
	identifierString := identifier.String()

	return s.acquire(
		identifierString,
		fmt.Sprintf("Realtime:user:%s:connections", identifier),
		connectionId.String(),
		maximumConnections,
	)
}

func (s *RealtimeLeaseStore) RefreshUserConnection(
	identifier uuid.UUID,
	connectionId uuid.UUID,
) (bool, error) {
	identifierString := identifier.String()

	return s.refresh(
		identifierString,
		fmt.Sprintf("Realtime:user:%s:connections", identifier),
		connectionId.String(),
	)
}

func (s *RealtimeLeaseStore) ReleaseUserConnection(
	identifier uuid.UUID,
	connectionId uuid.UUID,
) error {
	identifierString := identifier.String()

	return s.release(
		identifierString,
		fmt.Sprintf("Realtime:user:%s:connections", identifier),
		connectionId.String(),
	)
}

/* ============================== Block Pack Subscriber Methods ============================== */

func (s *RealtimeLeaseStore) AcquireBlockPackSubscriber(
	blockPackId uuid.UUID,
	member string,
	maximumSubscribers int,
) (bool, int64, error) {
	identifier := blockPackId.String()

	return s.acquire(
		identifier,
		fmt.Sprintf("Realtime:blockPack:%s:subscribers", blockPackId),
		member,
		maximumSubscribers,
	)
}

func (s *RealtimeLeaseStore) RefreshBlockPackSubscriber(
	blockPackId uuid.UUID,
	member string,
) (bool, error) {
	identifier := blockPackId.String()

	refreshed, err := s.refresh(
		identifier,
		fmt.Sprintf("Realtime:blockPack:%s:subscribers", blockPackId),
		member,
	)
	if err != nil || !refreshed {
		return refreshed, err
	}

	redisClient, err := s.getRedisClient(identifier)
	if err != nil {
		return false, err
	}

	if err := redisClient.PExpire(s.blockPackParticipantKey(blockPackId), constants.RealtimeLeaseTTL).Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (s *RealtimeLeaseStore) ReleaseBlockPackSubscriber(
	blockPackId uuid.UUID,
	member string,
) error {
	identifier := blockPackId.String()

	if err := s.release(
		identifier,
		fmt.Sprintf("Realtime:blockPack:%s:subscribers", blockPackId), member); err != nil {
		return err
	}

	redisClient, err := s.getRedisClient(identifier)
	if err != nil {
		return err
	}

	return redisClient.HDel(s.blockPackParticipantKey(blockPackId), member).Err()
}

/* ============================== Block Pack Participant Methods ============================== */

func (s *RealtimeLeaseStore) SetBlockPackParticipant(
	blockPackId uuid.UUID,
	member string,
	identifier uuid.UUID,
	channelPermission string,
) error {
	redisClient, err := s.getRedisClient(blockPackId.String())
	if err != nil {
		return err
	}

	payload, err := json.Marshal(RealtimeBlockPackParticipant{
		Member:            member,
		UserPublicId:      identifier.String(),
		ChannelPermission: channelPermission,
	})
	if err != nil {
		return err
	}

	pipeline := redisClient.TxPipeline()
	pipeline.HSet(s.blockPackParticipantKey(blockPackId), member, payload)
	pipeline.PExpire(s.blockPackParticipantKey(blockPackId), constants.RealtimeLeaseTTL)

	_, err = pipeline.Exec()

	return err
}

func (s *RealtimeLeaseStore) GetBlockPackParticipants(
	blockPackId uuid.UUID,
) ([]RealtimeBlockPackParticipant, error) {
	redisClient, err := s.getRedisClient(blockPackId.String())
	if err != nil {
		return nil, err
	}

	leaseKey := fmt.Sprintf("Realtime:blockPack:%s:subscribers", blockPackId)
	now := time.Now().UnixMilli()
	if err := redisClient.ZRemRangeByScore(leaseKey, "-inf", strconv.FormatInt(now, 10)).Err(); err != nil {
		return nil, err
	}

	members, err := redisClient.ZRangeByScore(
		leaseKey,
		redis.ZRangeBy{Min: strconv.FormatInt(now+1, 10), Max: "+inf"},
	).Result()
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return []RealtimeBlockPackParticipant{}, nil
	}

	values, err := redisClient.HMGet(s.blockPackParticipantKey(blockPackId), members...).Result()
	if err != nil {
		return nil, err
	}

	participants := make([]RealtimeBlockPackParticipant, 0, len(values))
	for _, value := range values {
		var payload []byte
		switch typedValue := value.(type) {
		case string:
			payload = []byte(typedValue)
		case []byte:
			payload = typedValue
		default:
			continue
		}

		var participant RealtimeBlockPackParticipant
		if err := json.Unmarshal(payload, &participant); err != nil {
			continue
		}

		participants = append(participants, participant)
	}

	return participants, nil
}
