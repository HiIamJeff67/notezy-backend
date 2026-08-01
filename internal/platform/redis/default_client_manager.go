package redis

import (
	"os"
	"strconv"

	configs "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
)

var defaultDatabaseNumber, _ = strconv.Atoi(os.Getenv("REDIS_INIT_DB"))

var DefaultClientManager = NewClientManager(configs.CacheManagerConfig{
	Host:     os.Getenv("REDIS_HOST"),
	Port:     os.Getenv("REDIS_PORT"),
	Password: os.Getenv("REDIS_PASSWORD"),
	DB:       defaultDatabaseNumber,
})

var ClientMap = DefaultClientManager.Clients()
