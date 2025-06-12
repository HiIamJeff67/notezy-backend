package caches

import (
	"context"
	"strconv"
	"sync"

	"github.com/go-redis/redis"

	exceptions "notezy-backend/app/exceptions"
	logs "notezy-backend/app/logs"
	shared "notezy-backend/app/shared"
	types "notezy-backend/app/shared/types"
)

var (
	RedisClientMap             map[int]*redis.Client                       = make(map[int]*redis.Client)
	RedisClientToConfig        map[*redis.Client]shared.CacheManagerConfig = make(map[*redis.Client]shared.CacheManagerConfig)
	PurposeToServerNumberRange                                             = map[shared.ValidCachePurpose]types.Range{
		shared.ValidCachePurpose_UserData:    UserDataRange,    // server number: 0 - 7 (included)
		shared.ValidCachePurpose_RecentPages: RecentPagesRange, // server number: 8 - 15 (included)
	}
	Ctx = context.Background()

	redisMapMutex sync.Mutex // since the map in go is not thread-safe, we need this mutex lock
)

func ConnectToRedis(config shared.CacheManagerConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		exceptions.Cache.FailedToConnectToServer(config.DB).WithError(err).Log().Panic()
	}

	redisMapMutex.Lock()
	defer redisMapMutex.Unlock()
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
		exceptions.Cache.ClientConfigDoesNotExist().Log()
		return false
	}

	if err := redisClient.Close(); err != nil {
		exceptions.Cache.FailedToDisconnectToServer(config.DB).WithError(err).Log()
		return false // since the server is just going to stop anyway, we don't need to panic here
	}

	redisMapMutex.Lock()
	defer redisMapMutex.Unlock()
	logs.FInfo("Deleting redis client server of %s into the RedisClientToConfig...", strconv.Itoa(config.DB))
	delete(RedisClientToConfig, redisClient)
	logs.FInfo("Deleting redis client server of %s into the RedisClientMap...", strconv.Itoa(config.DB))
	delete(RedisClientMap, config.DB)

	logs.FInfo("Redis client server of %s connected\n", strconv.Itoa(config.DB))

	return true
}

func ConnectToAllRedis() bool {
	var wg sync.WaitGroup                    // initialize the counter
	var resultCh chan bool = make(chan bool) // initialize the channel
	var totCount int = 0

	for _, serverRange := range PurposeToServerNumberRange {
		for i := serverRange.Start; i < serverRange.Start+serverRange.Size; i++ {
			totCount++
			wg.Add(1) // increase the counter by 1
			go func(dbIndex int) {
				defer wg.Done() // decrease the counter by 1 after this gorountine function returned
				currentConfig := shared.RedisCacheManagerConfigTemplate
				currentConfig.DB = dbIndex // modify the server number of the client
				res := ConnectToRedis(currentConfig)
				resultCh <- (res != nil)
			}(i)
		}
	}

	go func() {
		wg.Wait() // the wait group will stop here
		// once the counter is decreased back to 0, it will continue to close the resultCh
		close(resultCh)
	}()

	// the below part will end if the above gorountines are all finished
	var successCount int = 0
	for ok := range resultCh { // calculate the bool value in resultCh
		if ok {
			successCount++
		}
	}
	return successCount == totCount
}

func DisconnectToAllRedis() bool {
	var wg sync.WaitGroup
	var resultCh chan bool = make(chan bool)
	var totCount int = 0

	for _, serverRange := range PurposeToServerNumberRange {
		for i := serverRange.Start; i < serverRange.Start+serverRange.Size; i++ {
			totCount++
			wg.Add(1)
			go func(dbIndex int) {
				defer wg.Done()
				redisClient := RedisClientMap[dbIndex]
				ok := DisconnectToRedis(redisClient)
				resultCh <- !ok
			}(i)
		}
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var successCount int = 0
	for ok := range resultCh {
		if ok {
			successCount++
		}
	}
	return successCount == totCount
}
