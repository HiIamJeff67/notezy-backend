package config

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("GATEWAY_LISTEN_ADDRESS", "127.0.0.1:7777")
	t.Setenv("CORE_BASE_URL", "http://core:7778")
	t.Setenv("REALTIME_GATEWAY_BASE_URL", "http://realtime:7779")
	t.Setenv("CORE_CLIENT_TIMEOUT", "10s")
	t.Setenv("REALTIME_GATEWAY_CLIENT_TIMEOUT", "3s")
	t.Setenv("GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_START", "4")
	t.Setenv("GATEWAY_RATE_LIMIT_RECORD_CACHE_SERVER_SIZE", "4")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.CoreClientTimeout != 10*time.Second || config.RealtimeGatewayClientTimeout != 3*time.Second {
		t.Fatalf("LoadConfig() = %#v", config)
	}
}
