package config

import (
	"fmt"
	"os"
	"strings"

	platformpostgres "github.com/HiIamJeff67/notezy-backend/shared/platform/postgres"
)

func LoadPostgresConfig() (platformpostgres.Config, error) {
	config := platformpostgres.Config{
		Host:     strings.TrimSpace(os.Getenv("DB_HOST")),
		User:     strings.TrimSpace(os.Getenv("DB_USER")),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     strings.TrimSpace(os.Getenv("DB_NAME")),
		Port:     strings.TrimSpace(os.Getenv("DOCKER_DB_PORT")),
	}
	if config.Host == "" || config.User == "" || config.Password == "" || config.Name == "" || config.Port == "" {
		return platformpostgres.Config{}, fmt.Errorf("DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, and DOCKER_DB_PORT are required")
	}

	return config, nil
}
