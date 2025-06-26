package shared

/* ============================== Database Initialization ============================== */
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
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

/* ============================== API Initialization ============================== */
type CacheManagerConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

var (
	GinAddr                         = GetEnv("GIN_DOMAIN", "") + ":" + GetEnv("GIN_PORT", "7777")
	RedisCacheManagerConfigTemplate = CacheManagerConfig{
		Host:     GetEnv("REDIS_HOST", "notezy-redis"),
		Port:     GetEnv("REDIS_PORT", "6379"),
		Password: GetEnv("REDIS_PASSWORD", ""),
		DB:       GetIntEnv("REDIS_INIT_DB", 0),
	}
)
