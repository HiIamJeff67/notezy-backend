package database

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Setenv("DB_HOST", "database")
	t.Setenv("DB_USER", "notezy")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "notezy")
	t.Setenv("DOCKER_DB_PORT", "5432")

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if config.Host != "database" || config.Port != "5432" {
		t.Fatalf("LoadConfig() = %#v", config)
	}
}
