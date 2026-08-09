package configs

import (
	"fmt"
	"os"
	"strings"

	platformdatabase "github.com/HiIamJeff67/notezy-backend/shared/platform/database"
)

func LoadDatabaseConfig() (platformdatabase.Config, error) {
	config := platformdatabase.Config{
		Host:     strings.TrimSpace(os.Getenv("NOTIFICATION_DB_HOST")),
		User:     strings.TrimSpace(os.Getenv("NOTIFICATION_DB_USER")),
		Password: os.Getenv("NOTIFICATION_DB_PASSWORD"),
		Name:     strings.TrimSpace(os.Getenv("NOTIFICATION_DB_NAME")),
		Port:     strings.TrimSpace(os.Getenv("NOTIFICATION_DB_PORT")),
	}
	if config.Host == "" || config.User == "" || config.Password == "" || config.Name == "" || config.Port == "" {
		return platformdatabase.Config{}, fmt.Errorf("NOTIFICATION_DB_HOST, NOTIFICATION_DB_USER, NOTIFICATION_DB_PASSWORD, NOTIFICATION_DB_NAME, and NOTIFICATION_DB_PORT are required")
	}

	return config, nil
}
