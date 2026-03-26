package configs

import util "notezy-backend/app/util"

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string // the port inside the container, so please leave this as 5432 for PostgreSQL
}

var (
	PostgresDatabaseConfig = DatabaseConfig{
		Host:     util.GetEnv("DB_HOST", "notezy-db"),
		User:     util.GetEnv("DB_USER", "master"),
		Password: util.GetEnv("DB_PASSWORD", ""),
		DBName:   util.GetEnv("DB_NAME", "notezy-db"),
		Port:     util.GetEnv("DOCKER_DB_PORT", "5432"),
	}
)
