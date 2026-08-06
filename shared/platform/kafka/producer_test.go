package kafka

import (
	"testing"
	"time"
)

func TestNewProducerRejectsUnsupportedSASLMechanism(t *testing.T) {
	_, err := NewProducer(ClientConfig{
		ConnectionConfig: ConnectionConfig{
			Brokers: []string{
				"127.0.0.1:9094",
			},
			DialTimeout: time.Second,
			SASL: SASLConfig{
				Mechanism: "OAUTHBEARER",
				Username:  "test-user",
				Password:  "test-password",
			},
		},
		ClientId: "test-client",
	})
	if err == nil {
		t.Fatal("expected unsupported Kafka SASL mechanism to be rejected")
	}
}

func TestNewProducerRejectsPartialTLSClientCertificate(t *testing.T) {
	_, err := NewProducer(ClientConfig{
		ConnectionConfig: ConnectionConfig{
			Brokers: []string{
				"127.0.0.1:9094",
			},
			DialTimeout: time.Second,
			TLS: TLSConfig{
				Enabled:  true,
				CertFile: "client.crt",
			},
		},
		ClientId: "test-client",
	})
	if err == nil {
		t.Fatal("expected partial Kafka TLS client certificate to be rejected")
	}
}
