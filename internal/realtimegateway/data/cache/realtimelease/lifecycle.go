package realtimelease

import (
	"encoding/json"

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
