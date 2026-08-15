package config

import "testing"

func TestLoadPostgresConfig(t *testing.T) {
	t.Setenv("DB_HOST", "database")
	t.Setenv("DB_USER", "notezy")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "notezy")
	t.Setenv("DOCKER_DB_PORT", "5432")

	config, err := LoadPostgresConfig()
	if err != nil {
		t.Fatalf("LoadPostgresConfig() error = %v", err)
	}
	if config.Host != "database" || config.Port != "5432" {
		t.Fatalf("LoadPostgresConfig() = %#v", config)
	}
}
