package kafka

import (
	"context"
	"errors"
	"sync"

	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
)

var defaultProducerMutex sync.RWMutex
var DefaultProducer *Producer

func ConnectDefaultProducer(ctx context.Context) error {
	producer, err := NewProducer(config.Kafka())
	if err != nil {
		return err
	}

	defaultProducerMutex.Lock()
	previousProducer := DefaultProducer
	DefaultProducer = producer
	defaultProducerMutex.Unlock()
	if previousProducer != nil {
		previousProducer.Close()
	}

	return producer.Ping(ctx)
}

func CheckDefaultProducer(ctx context.Context) error {
	defaultProducerMutex.RLock()
	producer := DefaultProducer
	defaultProducerMutex.RUnlock()
	if producer == nil {
		return errors.New("Kafka producer is unavailable")
	}

	return producer.Ping(ctx)
}

func ProduceWithDefaultProducer(
	ctx context.Context,
	topic string,
	key string,
	value []byte,
) error {
	defaultProducerMutex.RLock()
	producer := DefaultProducer
	defaultProducerMutex.RUnlock()
	if producer == nil {
		return errors.New("Kafka producer is unavailable")
	}

	return producer.Produce(ctx, topic, key, value)
}

func CloseDefaultProducer() {
	defaultProducerMutex.Lock()
	producer := DefaultProducer
	DefaultProducer = nil
	defaultProducerMutex.Unlock()
	if producer != nil {
		producer.Close()
	}
}
