package redis

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host     string
	Port     string
	Password string
	Database int
}

type BackendServerName string

const (
	BackendServerName_EastAsia    BackendServerName = "EastAsia"
	BackendServerName_EastAmerica BackendServerName = "EastAmerica"
	BackendServerName_WestAmerica BackendServerName = "WestAmerica"
	BackendServerName_WestEurope  BackendServerName = "WestEurope"
)

var AllBackendServerNames = []BackendServerName{
	BackendServerName_EastAsia,
	BackendServerName_EastAmerica,
	BackendServerName_WestAmerica,
	BackendServerName_WestEurope,
}

func LoadConfig() (Config, error) {
	database, err := strconv.Atoi(strings.TrimSpace(os.Getenv("REDIS_INIT_DB")))
	if err != nil || database < 0 {
		return Config{}, fmt.Errorf("REDIS_INIT_DB must be a non-negative integer")
	}

	config := Config{
		Host:     strings.TrimSpace(os.Getenv("REDIS_HOST")),
		Port:     strings.TrimSpace(os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		Database: database,
	}
	if config.Host == "" || config.Port == "" {
		return Config{}, fmt.Errorf("REDIS_HOST and REDIS_PORT are required")
	}

	return config, nil
}
