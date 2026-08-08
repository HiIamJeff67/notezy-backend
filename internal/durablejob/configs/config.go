package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ListenAddress          string
	KafkaConsumer          KafkaConsumerConfig
	YjsMaintenanceStrategy YjsMaintenanceStrategyConfig
}

func LoadConfig() (Config, error) {
	config := Config{
		ListenAddress: strings.TrimSpace(os.Getenv("DURABLEJOB_LISTEN_ADDRESS")),
	}
	if config.ListenAddress == "" {
		return Config{}, fmt.Errorf("DURABLEJOB_LISTEN_ADDRESS is required")
	}

	var err error
	config.KafkaConsumer, err = loadKafkaConsumerConfig()
	if err != nil {
		return Config{}, err
	}
	config.YjsMaintenanceStrategy, err = loadYjsMaintenanceStrategyConfig()
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
