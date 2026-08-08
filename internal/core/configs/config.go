package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ListenAddress  string
	OAuthGoogle    OAuthGoogleConfig
	OutboxRelay    OutboxRelayConfig
	KafkaConsumer  KafkaConsumerConfig
	UserDataCache  UserDataCacheConfig
	StorageKeySalt string
}

func LoadConfig() (Config, error) {
	listenAddress := strings.TrimSpace(os.Getenv("CORE_LISTEN_ADDRESS"))
	if listenAddress == "" {
		return Config{}, fmt.Errorf("CORE_LISTEN_ADDRESS is required")
	}
	oauthGoogle, err := loadOAuthGoogleConfig()
	if err != nil {
		return Config{}, err
	}
	outboxRelay, err := loadOutboxRelayConfig()
	if err != nil {
		return Config{}, err
	}
	kafkaConsumer, err := loadKafkaConsumerConfig()
	if err != nil {
		return Config{}, err
	}
	storageKeySalt := os.Getenv("STORAGE_KEY_SALT")
	if storageKeySalt == "" {
		return Config{}, fmt.Errorf("STORAGE_KEY_SALT is required")
	}
	userDataCache, err := loadUserDataCacheConfig()
	if err != nil {
		return Config{}, err
	}
	return Config{
		ListenAddress:  listenAddress,
		OAuthGoogle:    oauthGoogle,
		OutboxRelay:    outboxRelay,
		KafkaConsumer:  kafkaConsumer,
		UserDataCache:  userDataCache,
		StorageKeySalt: storageKeySalt,
	}, nil
}
