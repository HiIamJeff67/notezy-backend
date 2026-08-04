package configs

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type KafkaConfig struct {
	Brokers       []string
	ClientId      string
	ConsumerGroup string
	DialTimeout   time.Duration
	TLS           KafkaTLSConfig
	SASL          KafkaSASLConfig
}

type KafkaTLSConfig struct {
	Enabled    bool
	CAFile     string
	CertFile   string
	KeyFile    string
	ServerName string
}

type KafkaSASLConfig struct {
	Mechanism string
	Username  string
	Password  string
}

func Kafka() KafkaConfig {
	configuredBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	brokers := make([]string, 0, len(configuredBrokers))
	for _, broker := range configuredBrokers {
		if trimmedBroker := strings.TrimSpace(broker); trimmedBroker != "" {
			brokers = append(brokers, trimmedBroker)
		}
	}
	if len(brokers) == 0 {
		brokers = []string{
			"127.0.0.1:9094",
		}
	}

	dialTimeoutSeconds, err := strconv.Atoi(os.Getenv("KAFKA_DIAL_TIMEOUT_SECONDS"))
	if err != nil || dialTimeoutSeconds <= 0 {
		dialTimeoutSeconds = 3
	}

	tlsEnabled, _ := strconv.ParseBool(os.Getenv("KAFKA_TLS_ENABLED"))
	clientId := os.Getenv("KAFKA_CLIENT_ID")
	if clientId == "" {
		clientId = "notezy-runtime"
	}
	consumerGroup := os.Getenv("KAFKA_CONSUMER_GROUP")
	if consumerGroup == "" {
		consumerGroup = "notezy-runtime"
	}

	return KafkaConfig{
		Brokers:       brokers,
		ClientId:      clientId,
		ConsumerGroup: consumerGroup,
		DialTimeout:   time.Duration(dialTimeoutSeconds) * time.Second,
		TLS: KafkaTLSConfig{
			Enabled:    tlsEnabled,
			CAFile:     os.Getenv("KAFKA_TLS_CA_FILE"),
			CertFile:   os.Getenv("KAFKA_TLS_CERT_FILE"),
			KeyFile:    os.Getenv("KAFKA_TLS_KEY_FILE"),
			ServerName: os.Getenv("KAFKA_TLS_SERVER_NAME"),
		},
		SASL: KafkaSASLConfig{
			Mechanism: strings.ToUpper(strings.TrimSpace(os.Getenv("KAFKA_SASL_MECHANISM"))),
			Username:  os.Getenv("KAFKA_SASL_USERNAME"),
			Password:  os.Getenv("KAFKA_SASL_PASSWORD"),
		},
	}
}
