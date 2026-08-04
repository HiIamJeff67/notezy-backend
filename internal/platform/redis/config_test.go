package redis

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Setenv("REDIS_HOST", "redis")
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_INIT_DB", "2")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.Database != 2 || config.Host != "redis" {
		t.Fatalf("LoadConfig() = %#v", config)
	}
}
