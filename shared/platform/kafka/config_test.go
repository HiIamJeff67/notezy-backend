package kafka

import (
	"testing"
	"time"
)

func TestLoadConnectionConfig(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", "kafka-1:9092, kafka-2:9092")
	t.Setenv("KAFKA_DIAL_TIMEOUT", "3s")
	t.Setenv("KAFKA_TLS_ENABLED", "false")

	config, err := LoadConnectionConfig()
	if err != nil {
		t.Fatalf("LoadConnectionConfig() error = %v", err)
	}
	if len(config.Brokers) != 2 || config.DialTimeout != 3*time.Second {
		t.Fatalf("LoadConnectionConfig() = %#v", config)
	}
}
