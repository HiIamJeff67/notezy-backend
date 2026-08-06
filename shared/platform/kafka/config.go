package kafka

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConnectionConfig struct {
	Brokers     []string
	DialTimeout time.Duration
	TLS         TLSConfig
	SASL        SASLConfig
}

type TLSConfig struct {
	Enabled    bool
	CAFile     string
	CertFile   string
	KeyFile    string
	ServerName string
}

type SASLConfig struct {
	Mechanism string
	Username  string
	Password  string
}

type ClientConfig struct {
	ConnectionConfig
	ClientId string
}

type ConsumerConfig struct {
	ClientConfig
	ConsumerGroup       string
	MaximumAttempts     int
	InitialRetryBackoff time.Duration
	MaximumRetryBackoff time.Duration
	MaximumPollRecords  int
}

func LoadConnectionConfig() (ConnectionConfig, error) {
	configuredBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	brokers := make([]string, 0, len(configuredBrokers))
	for _, broker := range configuredBrokers {
		if trimmedBroker := strings.TrimSpace(broker); trimmedBroker != "" {
			brokers = append(brokers, trimmedBroker)
		}
	}
	if len(brokers) == 0 {
		return ConnectionConfig{}, fmt.Errorf("KAFKA_BROKERS is required")
	}

	dialTimeout, err := time.ParseDuration(strings.TrimSpace(os.Getenv("KAFKA_DIAL_TIMEOUT")))
	if err != nil || dialTimeout <= 0 {
		return ConnectionConfig{}, fmt.Errorf("KAFKA_DIAL_TIMEOUT must be a positive Go duration")
	}

	tlsEnabled := false
	if value := strings.TrimSpace(os.Getenv("KAFKA_TLS_ENABLED")); value != "" {
		tlsEnabled, err = strconv.ParseBool(value)
	}
	if err != nil {
		return ConnectionConfig{}, fmt.Errorf("KAFKA_TLS_ENABLED must be a boolean")
	}

	return ConnectionConfig{
		Brokers:     brokers,
		DialTimeout: dialTimeout,
		TLS: TLSConfig{
			Enabled:    tlsEnabled,
			CAFile:     strings.TrimSpace(os.Getenv("KAFKA_TLS_CA_FILE")),
			CertFile:   strings.TrimSpace(os.Getenv("KAFKA_TLS_CERT_FILE")),
			KeyFile:    strings.TrimSpace(os.Getenv("KAFKA_TLS_KEY_FILE")),
			ServerName: strings.TrimSpace(os.Getenv("KAFKA_TLS_SERVER_NAME")),
		},
		SASL: SASLConfig{
			Mechanism: strings.ToUpper(strings.TrimSpace(os.Getenv("KAFKA_SASL_MECHANISM"))),
			Username:  strings.TrimSpace(os.Getenv("KAFKA_SASL_USERNAME")),
			Password:  os.Getenv("KAFKA_SASL_PASSWORD"),
		},
	}, nil
}
