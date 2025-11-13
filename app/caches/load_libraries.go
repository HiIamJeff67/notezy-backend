package caches

import (
	"fmt"

	redisfunctionlibraries "notezy-backend/app/caches/libraries"
	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
)

func FlushRedisFunctionsLibraries() *exceptions.Exception {
	for serverName, serverNumber := range BackendServerNameToRateLimitRedisIndex {
		redisClient, exist := RedisClientMap[serverNumber]
		if !exist {
			continue
		}

		redisClient.Do("FUNCTION", "FLUSH")
		logs.FDebug("Flushed all the functions across all libraries in server %s of %d", serverName, serverNumber)
	}

	return nil
}

func ReloadRateLimitRecordRedisFunctionsLibraries() *exceptions.Exception {
	for serverName, serverNumber := range BackendServerNameToRateLimitRedisIndex {
		redisClient, exist := RedisClientMap[serverNumber]
		if !exist {
			continue
		}

		if err := redisClient.Do("FUNCTION", "LOAD", "REPLACE", redisfunctionlibraries.RateLimitRecordRedisFunctionsLibraryContent).Err(); err != nil {
			return exceptions.Cache.FailedToLoadRedisFunctions().
				WithDetails(fmt.Sprintf("Failed to load functions from lua scripts in server %s of %d", serverName, serverNumber)).
				WithError(err)
		}

		logs.FInfo("Reloaded all the functions in library of %s from lua scripts in server %s of %d",
			redisfunctionlibraries.RateLimitRecordRedisFunctionsLibrary,
			serverName,
			serverNumber,
		)
	}

	return nil
}
