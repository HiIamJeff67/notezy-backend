package realtimelease

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
)

type BlockPackChannelRevocation struct {
	EventId            uuid.UUID                                          `json:"eventId"`
	BlockPackId        uuid.UUID                                          `json:"blockPackId"`
	TargetUserPublicId *uuid.UUID                                         `json:"targetUserPublicId,omitempty"`
	Reason             coreeventscontract.BlockPackAccessRevocationReason `json:"reason"`
}

type UserSessionRevocation struct {
	EventId      uuid.UUID `json:"eventId"`
	UserPublicId uuid.UUID `json:"userPublicId"`
}

type ResourceEvent struct {
	EventId            uuid.UUID  `json:"eventId"`
	EventType          string     `json:"eventType"`
	ResourceId         uuid.UUID  `json:"resourceId"`
	TargetUserPublicId *uuid.UUID `json:"targetUserPublicId,omitempty"`
	Change             string     `json:"change"`
	Permission         string     `json:"permission,omitempty"`
}

func (s *RealtimeLeaseCacheClient) MarkLifecycleEventProcessed(eventId uuid.UUID) (bool, error) {
	if eventId == uuid.Nil {
		return false, fmt.Errorf("realtime lifecycle event ID is required")
	}

	redisClient, err := s.getRedisClient("lifecycle-events")
	if err != nil {
		return false, err
	}

	return redisClient.SetNX(
		fmt.Sprintf("Realtime:lifecycle:event:%s", eventId),
		"1",
		7*24*time.Hour,
	).Result()
}

func (s *RealtimeLeaseCacheClient) ReleaseLifecycleEvent(eventId uuid.UUID) error {
	redisClient, err := s.getRedisClient("lifecycle-events")
	if err != nil {
		return err
	}

	return redisClient.Del(fmt.Sprintf("Realtime:lifecycle:event:%s", eventId)).Err()
}

func (s *RealtimeLeaseCacheClient) PublishBlockPackChannelRevocation(
	revocation BlockPackChannelRevocation,
) error {
	redisClient, err := s.getRedisClient(revocation.BlockPackId.String())
	if err != nil {
		return err
	}

	payload, err := json.Marshal(revocation)
	if err != nil {
		return err
	}

	return redisClient.Publish(s.blockPackChannelRevocationKey(), payload).Err()
}

func (s *RealtimeLeaseCacheClient) SubscribeBlockPackChannelRevocations(
	handler func(BlockPackChannelRevocation),
) (func(), error) {
	redisClient, err := s.getRedisClient(s.blockPackChannelRevocationKey())
	if err != nil {
		return nil, err
	}

	pubsub := redisClient.Subscribe(s.blockPackChannelRevocationKey())

	go func() {
		for message := range pubsub.Channel() {
			var revocation BlockPackChannelRevocation
			if err := json.Unmarshal([]byte(message.Payload), &revocation); err != nil {
				continue
			}

			handler(revocation)
		}
	}()

	return func() {
		_ = pubsub.Close()
	}, nil
}

func (s *RealtimeLeaseCacheClient) PublishUserSessionRevocation(
	revocation UserSessionRevocation,
) error {
	redisClient, err := s.getRedisClient(revocation.UserPublicId.String())
	if err != nil {
		return err
	}

	payload, err := json.Marshal(revocation)
	if err != nil {
		return err
	}

	return redisClient.Publish(s.userSessionRevocationKey(), payload).Err()
}

func (s *RealtimeLeaseCacheClient) SubscribeUserSessionRevocations(
	handler func(UserSessionRevocation),
) (func(), error) {
	redisClient, err := s.getRedisClient(s.userSessionRevocationKey())
	if err != nil {
		return nil, err
	}

	pubsub := redisClient.Subscribe(s.userSessionRevocationKey())

	go func() {
		for message := range pubsub.Channel() {
			var revocation UserSessionRevocation
			if err := json.Unmarshal([]byte(message.Payload), &revocation); err != nil {
				continue
			}

			handler(revocation)
		}
	}()

	return func() {
		_ = pubsub.Close()
	}, nil
}

func (s *RealtimeLeaseCacheClient) PublishResourceEvent(event ResourceEvent) error {
	redisClient, err := s.getRedisClient(s.resourceEventKey())
	if err != nil {
		return err
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return redisClient.Publish(s.resourceEventKey(), payload).Err()
}

func (s *RealtimeLeaseCacheClient) SubscribeResourceEvents(
	handler func(ResourceEvent),
) (func(), error) {
	redisClient, err := s.getRedisClient(s.resourceEventKey())
	if err != nil {
		return nil, err
	}

	pubsub := redisClient.Subscribe(s.resourceEventKey())

	go func() {
		for message := range pubsub.Channel() {
			var event ResourceEvent
			if err := json.Unmarshal([]byte(message.Payload), &event); err != nil {
				continue
			}

			handler(event)
		}
	}()

	return func() {
		_ = pubsub.Close()
	}, nil
}

func (s *RealtimeLeaseCacheClient) resourceEventKey() string {
	return "Realtime:resource:events"
}
