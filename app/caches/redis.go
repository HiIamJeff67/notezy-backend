package caches

import (
	"context"
	logs "go-gorm-api/app/logs"
	"go-gorm-api/global"
	"strconv"

	"github.com/go-redis/redis"
)

var (
	RedisClientMap map[int]*redis.Client
	RedisClientToConfig map[*redis.Client]global.CacheManagerConfig
	PurposeToServerNumberRange = map[global.ValidCachePurpose]global.Range{
		global.ValidCachePurpose_UserData: UserDataRange, 
		global.ValidCachePurpose_RecentPages: RecentPagesRange,
	}
	Ctx = context.Background()
)

func ConnectToRedis(config global.CacheManagerConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Host + ":" + config.Port, 
		Password: config.Password,
		DB: config.DB, 
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		logs.FError("Error connecting to the redis client server of %s\n", strconv.Itoa(config.DB))
		panic("Connecting to redis client error : " + err.Error())
	}

	if _, ok := RedisClientToConfig[redisClient]; !ok {
		logs.FInfo("Storing redis client server of %s into the RedisClientToConfig...", strconv.Itoa(config.DB))
		RedisClientToConfig[redisClient] = config
	}
	if _, ok := RedisClientMap[config.DB]; !ok {
		logs.FInfo("Storing redis client server of %s into the RedisClientMap...", strconv.Itoa(config.DB))
		RedisClientMap[config.DB] = redisClient
	}

	logs.FInfo("Redis client server of %s connected\n", strconv.Itoa(config.DB))

	return redisClient
}

func DisconnectToRedis(redisClient *redis.Client) bool {
	config, ok := RedisClientToConfig[redisClient]
	if !ok {
		logs.FError("Failed to get the connection of the given redis")
		return false
	}

	if err := redisClient.Close(); err != nil {
		logs.FError("Error connecting to the redis client server of %s\n", strconv.Itoa(config.DB))
		return false
	}
	
	logs.FInfo("Extracting redis client server of %s into the RedisClientToConfig...", strconv.Itoa(config.DB))
	delete(RedisClientToConfig, redisClient)
	logs.FInfo("Extracting redis client server of %s into the RedisClientMap...", strconv.Itoa(config.DB))
	delete(RedisClientMap, config.DB)

	logs.FInfo("Redis client server of %s connected\n", strconv.Itoa(config.DB))

	return true
}

func ConnectToAllRedis() bool {
	var count int = 0; var totCount int = 0;
	for _, serverRange := range PurposeToServerNumberRange {
		for i := serverRange.Start; i < serverRange.Size; i++ {
			currentConfig := global.RedisCacheManagerConfigTemplate
			currentConfig.DB = i
			if res := ConnectToRedis(currentConfig); res != nil { count++ }
		}
		totCount += serverRange.Size
	}
	return count == totCount
}

func DisconnectToAllRedis() bool {
	var count int = 0; var totCount int = 0;
	for _, serverRange := range PurposeToServerNumberRange {
		for i := serverRange.Start; i < serverRange.Size; i++ {
			redisClient := RedisClientMap[i]
			if ok := DisconnectToRedis(redisClient); !ok { count++ }
		}
		totCount += serverRange.Size
	}
	return count == totCount
}