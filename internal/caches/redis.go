package caches

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"

	redislibraries "github.com/HiIamJeff67/notezy-backend/internal/caches/libraries"
	exceptions "github.com/HiIamJeff67/notezy-backend/internal/exceptions"
	logs "github.com/HiIamJeff67/notezy-backend/internal/platform/observability/logs"
	platformredis "github.com/HiIamJeff67/notezy-backend/internal/platform/redis"
	constants "github.com/HiIamJeff67/notezy-backend/internal/shared/constants"
	types "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

var (
	RedisClientMap             map[int]*redis.Client = platformredis.ClientMap
	UserDataStore                                    = NewUserDataCacheStore(RedisClientMap)
	RateLimitRecordStore                             = NewRateLimitRecordCacheStore(RedisClientMap)
	PurposeToServerNumberRange                       = map[types.ValidCachePurpose]types.Range[int, int]{
		types.ValidCachePurpose_UserData:   UserDataStore.Range,        // server number: 0 - 3 (included)
		types.ValidCachePurpose_RateLimite: RateLimitRecordStore.Range, // server number: 4 - 7 (included)
		types.ValidCachePurpose_Realtime:   types.Range[int, int]{Start: constants.RealtimeRedisServerNumber, Size: 1},
	}
)

func FlushCacheLibraries() *exceptions.Exception {
	for serverName, serverNumber := range RateLimitRecordStore.backendServerNameToRedisIndex {
		redisClient, exist := RedisClientMap[serverNumber]
		if !exist {
			continue
		}

		redisClient.Do("FUNCTION", "FLUSH")
		logs.NotezyLogger.Debug(context.Background(), fmt.Sprintf("Flushed all the functions across all libraries in server %s of %d", serverName, serverNumber))
	}

	return nil
}

func LoadRateLimitRecordCacheLibraries() *exceptions.Exception {
	for serverName, serverNumber := range RateLimitRecordStore.backendServerNameToRedisIndex {
		redisClient, exist := RedisClientMap[serverNumber]
		if !exist {
			continue
		}

		if err := redisClient.Do("FUNCTION", "LOAD", "REPLACE", redislibraries.RateLimitRecordLibraryContent).Err(); err != nil {
			return exceptions.New(
				"RedisFunctionLoadFailed",
				"Cache",
				"LoadRateLimitRecordCacheLibraries",
				"Failed to load rate limit cache functions",
				http.StatusInternalServerError,
				true,
			).WithOrigin(fmt.Errorf(
				"failed to load functions from Lua scripts in server %s of %d: %w",
				serverName,
				serverNumber,
				err,
			))
		}

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Reloaded all the functions in library of %s from lua scripts in server %s of %d", redislibraries.RateLimitRecordLibrary, serverName, serverNumber))
	}

	return nil
}

func LoadUserQuotaCacheLibraries() *exceptions.Exception {
	for serverNumber := UserDataStore.Range.Start; serverNumber < UserDataStore.Range.Start+UserDataStore.Range.Size; serverNumber++ {
		redisClient, exist := RedisClientMap[serverNumber]
		if !exist {
			continue
		}

		if err := redisClient.Do("FUNCTION", "LOAD", "REPLACE", redislibraries.UserQuotaLibraryContent).Err(); err != nil {
			return exceptions.New(
				"RedisFunctionLoadFailed",
				"Cache",
				"LoadUserQuotaCacheLibraries",
				"Failed to load user quota cache functions",
				http.StatusInternalServerError,
				true,
			).WithOrigin(fmt.Errorf(
				"failed to load functions from Lua scripts in server number of %d: %w",
				serverNumber,
				err,
			))
		}

		logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Reloaded all the functions in library of %s from lua scripts in server number of %d", redislibraries.UserQuotaLibrary, serverNumber))
	}

	return nil
}
