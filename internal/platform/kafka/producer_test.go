package kafka

import (
	"testing"
	"time"

	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
)

func TestNewProducerRejectsUnsupportedSASLMechanism(t *testing.T) {
	_, err := NewProducer(config.KafkaConfig{
		Brokers: []string{
			"127.0.0.1:9094",
		},
		ClientId:      "test-client",
		ConsumerGroup: "test-group",
		DialTimeout:   time.Second,
		SASL: config.KafkaSASLConfig{
			Mechanism: "OAUTHBEARER",
			Username:  "test-user",
			Password:  "test-password",
		},
	})
	if err == nil {
		t.Fatal("expected unsupported Kafka SASL mechanism to be rejected")
	}
}

func TestNewProducerRejectsPartialTLSClientCertificate(t *testing.T) {
	_, err := NewProducer(config.KafkaConfig{
		Brokers: []string{
			"127.0.0.1:9094",
		},
		ClientId:      "test-client",
		ConsumerGroup: "test-group",
		DialTimeout:   time.Second,
		TLS: config.KafkaTLSConfig{
			Enabled:  true,
			CertFile: "client.crt",
		},
	})
	if err == nil {
		t.Fatal("expected partial Kafka TLS client certificate to be rejected")
	}
}
