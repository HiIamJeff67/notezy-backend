package lib

import (
	"sync"
	"time"

	"github.com/google/uuid"
	rate "golang.org/x/time/rate"

	caches "notezy-backend/app/caches"
	logs "notezy-backend/app/logs"
	"notezy-backend/shared/constants"
	types "notezy-backend/shared/types"
)

type HybridRateLimitTask struct {
	NumOfChangingTokens int32 `json:"numOfChangingTokens"`
	IsAccumulated       bool  `json:"isAccumulated"`
	Retries             int   `json:"retries"`
	MaxRetries          int   `json:"maxRetries"`
}

type HybridRateLimiter struct {
	*rate.Limiter
	LimiterMutex   sync.RWMutex
	UserLimit      int32
	WindowDuration time.Duration

	pendingTasks      map[string]HybridRateLimitTask
	pendingTasksMutex sync.Mutex
	syncInterval      time.Duration
	syncTicker        *time.Ticker
	stopChan          chan struct{}

	BackendServerName   types.BackendServerName
	IsAuthorizedLimiter bool
}

func NewHybridRateLimiter(
	rateLimit rate.Limit,
	burst int,
	userLimit int32,
	windowDuration time.Duration,
	backendServerName types.BackendServerName,
	isAuthorizedLimiter bool,
) *HybridRateLimiter {
	syncInterval := windowDuration / constants.SynchronizationToWindowDurationRatio
	syncInterval = max(constants.MinSynchornizationInterval, syncInterval)

	hrl := &HybridRateLimiter{
		Limiter:             rate.NewLimiter(rateLimit, burst),
		UserLimit:           userLimit,
		WindowDuration:      windowDuration,
		pendingTasks:        make(map[string]HybridRateLimitTask, 0),
		syncInterval:        syncInterval,
		syncTicker:          time.NewTicker(syncInterval),
		stopChan:            make(chan struct{}),
		BackendServerName:   backendServerName,
		IsAuthorizedLimiter: isAuthorizedLimiter,
	}

	// initially calling syncLoop()
	go hrl.syncLoop()

	return hrl
}

func (hrl *HybridRateLimiter) reappendTasks(failedTasks map[string]HybridRateLimitTask) {
	hrl.pendingTasksMutex.Lock()
	defer hrl.pendingTasksMutex.Unlock()

	for key, task := range failedTasks {
		if task.Retries < task.MaxRetries {
			hrl.pendingTasks[key] = HybridRateLimitTask{
				NumOfChangingTokens: task.NumOfChangingTokens,
				IsAccumulated:       task.IsAccumulated,
				Retries:             task.Retries + 1,
				MaxRetries:          task.MaxRetries,
			}
		} else {
			logs.FWarn("Dropping task for key %s after %d retries", key, task.MaxRetries)
		}
	}
}

func (hrl *HybridRateLimiter) batchSync() {
	hrl.pendingTasksMutex.Lock()
	if len(hrl.pendingTasks) == 0 {
		hrl.pendingTasksMutex.Unlock()
		return
	}

	syncData := make(map[string]HybridRateLimitTask)
	for key, task := range hrl.pendingTasks {
		syncData[key] = task
	}
	hrl.pendingTasks = make(map[string]HybridRateLimitTask)
	hrl.pendingTasksMutex.Unlock()

	if hrl.IsAuthorizedLimiter {
		userDtos := make([]struct {
			UserId         uuid.UUID                                 `json:"userId"`
			SynchronizeDto caches.SynchronizeRateLimitRecordCacheDto `json:"synchronizeDto"`
		}, 0, len(syncData))

		for userIdStr, task := range syncData {
			userId, err := uuid.Parse(userIdStr)
			if err != nil {
				logs.FError("Failed to parse user ID %s: %v", userIdStr, err)
				continue
			}

			userDtos = append(userDtos, struct {
				UserId         uuid.UUID                                 `json:"userId"`
				SynchronizeDto caches.SynchronizeRateLimitRecordCacheDto `json:"synchronizeDto"`
			}{
				UserId: userId,
				SynchronizeDto: caches.SynchronizeRateLimitRecordCacheDto{
					NumOfChangingTokens: task.NumOfChangingTokens,
					IsAccumulated:       task.IsAccumulated,
				},
			})
		}

		if err := caches.BatchSynchronizeRateLimitRecordCachesByUserIds(userDtos, hrl.BackendServerName); err != nil {
			logs.FError("Failed to batch sync user rate limits to Redis: %v", err)
			hrl.reappendTasks(syncData)
		} else if len(userDtos) > 0 {
			logs.FInfo("Batch synced %d user rate limits to Redis", len(userDtos))
		}
	} else {
		clientDtos := make([]struct {
			Fingerprint    string                                    `json:"fingerprint"`
			SynchronizeDto caches.SynchronizeRateLimitRecordCacheDto `json:"synchronizeDto"`
		}, 0, len(syncData))

		for fingerprint, task := range syncData {
			clientDtos = append(clientDtos, struct {
				Fingerprint    string                                    `json:"fingerprint"`
				SynchronizeDto caches.SynchronizeRateLimitRecordCacheDto `json:"synchronizeDto"`
			}{
				Fingerprint: fingerprint,
				SynchronizeDto: caches.SynchronizeRateLimitRecordCacheDto{
					NumOfChangingTokens: task.NumOfChangingTokens,
					IsAccumulated:       task.IsAccumulated,
				},
			})
		}

		if err := caches.BatchSynchronizeRateLimitRecordCachesByFingerprints(clientDtos, hrl.BackendServerName); err != nil {
			logs.FError("Failed to batch sync client IP rate limits to Redis: %v", err)
			hrl.reappendTasks(syncData)
		} else if len(clientDtos) > 0 {
			logs.FInfo("Batch synced %d client IP rate limits to Redis", len(clientDtos))
		}
	}
}

func (hrl *HybridRateLimiter) appendPendingTask(key string, tokens int32) {
	hrl.pendingTasksMutex.Lock()
	defer hrl.pendingTasksMutex.Unlock()

	if existingTask, exists := hrl.pendingTasks[key]; exists {
		hrl.pendingTasks[key] = HybridRateLimitTask{
			NumOfChangingTokens: existingTask.NumOfChangingTokens + tokens,
			IsAccumulated:       true,
			Retries:             existingTask.Retries,
			MaxRetries:          3,
		}
	} else {
		hrl.pendingTasks[key] = HybridRateLimitTask{
			NumOfChangingTokens: tokens,
			IsAccumulated:       true,
			Retries:             0,
			MaxRetries:          3,
		}
	}
}

func (hrl *HybridRateLimiter) syncLoop() {
	for {
		select {
		case <-hrl.syncTicker.C:
			hrl.batchSync()
		case <-hrl.stopChan:
			hrl.batchSync()
			return
		}
	}
}

/* ============================== Operations for Unauthorized Middleware (Client IP based) ============================== */

func (hrl *HybridRateLimiter) checkBucketLimitByFingerprint(fingerprint string, n int32) int32 {
	var totalTokensUsed int32 = 0

	for _, backendServerName := range types.AllBackendServerNames {
		rateLimitRecordCache, exception := caches.GetRateLimitRecordCacheByFingerprint(fingerprint, backendServerName)
		if exception != nil {
			continue
		}

		now := time.Now()
		if now.Sub(rateLimitRecordCache.WindowStartTime) >= rateLimitRecordCache.WindowDuration {
			continue
		}

		totalTokensUsed += rateLimitRecordCache.NumOfTokens
	}

	return max(0, hrl.UserLimit-totalTokensUsed-n)
}

func (hrl *HybridRateLimiter) AllowByFingerprint(fingerprint string) (bool, int32) {
	return hrl.AllowNByFingerprint(fingerprint, time.Now(), 1)
}

func (hrl *HybridRateLimiter) AllowNByFingerprint(fingerprint string, now time.Time, n int) (bool, int32) {
	hrl.LimiterMutex.RLock()
	defer hrl.LimiterMutex.RUnlock()

	// 1. Use the Limiter from the rate utility for fast checking
	if !hrl.Limiter.AllowN(now, n) {
		logs.FInfo("Request blocked by local rate limiter for client IP: %s, requested: %d", fingerprint, n)
		return false, 0
	}

	// 2. Use the rate limit record in redis cache to check if the request from the same source has exceeded some certain count
	remaining := hrl.checkBucketLimitByFingerprint(fingerprint, int32(n))
	if remaining < 0 {
		logs.FInfo("Request blocked by global rate limiter for client IP: %s, requested: %d", fingerprint, n)
		return false, 0
	}

	// 3. Record the tokens that being used, and batch synchronize to the redis
	hrl.appendPendingTask(fingerprint, int32(n))

	return true, remaining
}

/* ============================== Operations for Authorized Middleware (User ID based) ============================== */

func (hrl *HybridRateLimiter) checkBucketLimitByUserId(userId uuid.UUID, fingerprint string, n int32) int32 {
	var totalTokensUsed int32 = 0

	for _, backendServerName := range types.AllBackendServerNames {
		rateLimitRecordCache, exception := caches.GetRateLimitRecordCacheByUserId(userId, fingerprint, backendServerName)
		if exception != nil {
			continue
		}

		now := time.Now()
		if now.Sub(rateLimitRecordCache.WindowStartTime) >= rateLimitRecordCache.WindowDuration {
			continue
		}

		totalTokensUsed += rateLimitRecordCache.NumOfTokens
	}

	logs.Info("The current tokens used by the user: ", totalTokensUsed)

	return max(0, hrl.UserLimit-totalTokensUsed-n)
}

func (hrl *HybridRateLimiter) AllowByUserId(userId uuid.UUID, fingerprint string) (bool, int32) {
	return hrl.AllowNByUserId(userId, fingerprint, time.Now(), 1)
}

func (hrl *HybridRateLimiter) AllowNByUserId(userId uuid.UUID, fingerprint string, now time.Time, n int) (bool, int32) {
	hrl.LimiterMutex.RLock()
	defer hrl.LimiterMutex.RUnlock()

	// 1. Use the Limiter from the rate utility for fast checking
	if !hrl.Limiter.AllowN(now, n) {
		logs.FInfo("Request blocked by local rate limiter for user ID: %s, requested: %d", userId.String(), n)
		return false, 0
	}

	// 2. Use the rate limit record in redis cache to check if the request from the same source has exceeded some certain count
	remaining := hrl.checkBucketLimitByUserId(userId, fingerprint, int32(n))
	if remaining < 0 {
		logs.FInfo("Request blocked by global rate limiter for user ID: %s, requested: %d", userId.String(), n)
		return false, 0
	}

	// 3. Record the tokens that being used, and batch synchronize to the redis
	hrl.appendPendingTask(userId.String(), int32(n))

	return true, remaining
}

/* ============================== Utility Methods ============================== */

func (hrl *HybridRateLimiter) Allow(key string) (bool, int32) {
	if hrl.IsAuthorizedLimiter {
		userId, err := uuid.Parse(key)
		if err != nil {
			logs.FError("Invalid user ID format: %s", key)
			return false, 0
		}
		return hrl.AllowByUserId(userId, "")
	} else {
		return hrl.AllowByFingerprint(key)
	}
}

func (hrl *HybridRateLimiter) AllowN(key string, now time.Time, n int) (bool, int32) {
	if hrl.IsAuthorizedLimiter {
		userId, err := uuid.Parse(key)
		if err != nil {
			logs.FError("Invalid user ID format: %s", key)
			return false, 0
		}
		return hrl.AllowNByUserId(userId, "", now, n)
	} else {
		return hrl.AllowNByFingerprint(key, now, n)
	}
}

func (hrl *HybridRateLimiter) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"userLimit":           hrl.UserLimit,
		"windowDuration":      hrl.WindowDuration,
		"backendServerName":   hrl.BackendServerName,
		"isAuthorizedLimiter": hrl.IsAuthorizedLimiter,
	}
}

func (hrl *HybridRateLimiter) GetDetailStatus() map[string]interface{} {
	hrl.pendingTasksMutex.Lock()
	defer hrl.pendingTasksMutex.Unlock()
	return map[string]interface{}{
		"userLimit":           hrl.UserLimit,
		"windowDuration":      hrl.WindowDuration,
		"isPending":           len(hrl.pendingTasks) > 0,
		"backendServerName":   hrl.BackendServerName,
		"isAuthorizedLimiter": hrl.IsAuthorizedLimiter,
	}
}

func (hrl *HybridRateLimiter) Stop() {
	close(hrl.stopChan)
	hrl.syncTicker.Stop()
}
