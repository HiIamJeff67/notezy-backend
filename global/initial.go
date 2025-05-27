package global

import (
	"os"
	"strconv"
)

/* ============================== Database Initialization ============================== */
type DatabaseConfig struct {
	Host string 
	User string
	Password string
	DBName string
	Port string
}

var (
	PostgresDatabaseConfig = DatabaseConfig{
		Host:     GetEnv("DB_HOST", "notezy-db"),
		User:     GetEnv("DB_USER", "master"),
		Password: GetEnv("DB_PASSWORD", ""),
		DBName:   GetEnv("DB_NAME", "notezy-db"),
		Port:     GetEnv("DOCKER_DB_PORT", "5432"), // we should use the port hosted by docker
	}
)
/* ============================== Database Initialization ============================== */

/* ============================== API Initialization ============================== */
type CacheManagerConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

var (
	GinAddr = GetEnv("GIN_DOMAIN", "") + ":" + GetEnv("GIN_PORT", "7777")
	RedisCacheManagerConfigTemplate = CacheManagerConfig{
		Host: GetEnv("REDIS_HOST", "notezy-redis"), 
		Port: GetEnv("REDIS_PORT", "6379"), 
		Password: GetEnv("REDIS_PASSWORD", ""), 
		DB: GetIntEnv("REDIS_INIT_DB", 0), 
	}
)
/* ============================== API Initialization ============================== */

/* ============================== Temporary Environment Variables Fetcher ============================== */
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetIntEnv(key string, fallback int) int {
	if valueStr, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fallback
		}
		return value
	}
	return fallback
}
/* ============================== Temporary Environment Variables Fetcher ============================== */