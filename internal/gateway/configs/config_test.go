package config

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("GATEWAY_LISTEN_ADDRESS", "127.0.0.1:7777")
	t.Setenv("CORE_BASE_URL", "http://core:7778")
	t.Setenv("CORE_CLIENT_TIMEOUT", "10s")
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.CoreAdapterTimeout != 10*time.Second {
		t.Fatalf("LoadConfig() = %#v", config)
	}
}
