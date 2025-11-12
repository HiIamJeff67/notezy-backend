package redisfunctions

import (
	"os"

	"notezy-backend/app/caches"
	exceptions "notezy-backend/app/exceptions"
	"notezy-backend/app/logs"
)

func LoadRateLimitRecordRedisFunctions() error {
	functionScriptsContent, err := os.ReadFile("rate_limit_record_functions.lua")
	if err != nil {
		return exceptions.Cache.FileExceptionDomain.CannotOpenFiles().Error
	}

	for serverName, serverNumber := range caches.BackendServerNameToRateLimitRedisIndex {
		redisClient, exist := caches.RedisClientMap[serverNumber]
		if !exist {
			continue
		}

		result := redisClient.Do("FUNCTION", "LOAD", string(functionScriptsContent))
		if err := result.Err(); err != nil {
			logs.FError("Failed to load functions from lua scripts in server %s of %d", serverName, serverNumber)
			return err
		}

		logs.FInfo("Successfully load functions from lua scripts in server %s of %d", serverName, serverNumber)
	}

	return nil
}
