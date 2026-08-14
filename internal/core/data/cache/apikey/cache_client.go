package apikey

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	redisclient "github.com/go-redis/redis"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"
	platformredis "github.com/HiIamJeff67/notezy-backend/shared/platform/redis"
)

const defaultCacheExpiresIn = 5 * time.Minute

type APIKeyCache struct {
	Id           uuid.UUID  `json:"id"`
	UserId       uuid.UUID  `json:"userId"`
	UserPublicId uuid.UUID  `json:"userPublicId"`
	ExpiresAt    *time.Time `json:"expiresAt"`
	RevokedAt    *time.Time `json:"revokedAt"`
}

type APIKeyCacheClient struct {
	clientSet *platformredis.ClientSet
	expiresIn time.Duration
}

func NewAPIKeyCacheClient(clientSet *platformredis.ClientSet, expiresIn ...time.Duration) *APIKeyCacheClient {
	ttl := defaultCacheExpiresIn
	if len(expiresIn) > 0 && expiresIn[0] > 0 {
		ttl = expiresIn[0]
	}
	return &APIKeyCacheClient{clientSet: clientSet, expiresIn: ttl}
}

func formatAPIKeyCacheKey(keyHash string) string {
	return fmt.Sprintf("APIKey:%s", keyHash)
}

func (c *APIKeyCacheClient) redisClient(keyHash string) (*redisclient.Client, *exceptions.Exception) {
	if c == nil || c.clientSet == nil {
		return nil, exceptions.New("CacheClientUnavailable", "Cache", "GetAPIKey", "API key cache is unavailable", http.StatusInternalServerError, true)
	}
	client, _, err := c.clientSet.ClientForKey(keyHash)
	if err != nil || client == nil {
		return nil, exceptions.New("CacheClientUnavailable", "Cache", "GetAPIKey", "API key cache is unavailable", http.StatusInternalServerError, true).WithOrigin(err)
	}
	return client, nil
}

func (c *APIKeyCacheClient) Get(keyHash string) (*APIKeyCache, *exceptions.Exception) {
	client, exception := c.redisClient(keyHash)
	if exception != nil {
		return nil, exception
	}
	value, err := client.Get(formatAPIKeyCacheKey(keyHash)).Result()
	if err == redisclient.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, exceptions.New("CacheReadFailed", "Cache", "GetAPIKey", "API key cache read failed", http.StatusInternalServerError, true).WithOrigin(err)
	}
	cache := &APIKeyCache{}
	if err := json.Unmarshal([]byte(value), cache); err != nil {
		return nil, exceptions.New("CacheDecodeFailed", "Cache", "GetAPIKey", "API key cache value is invalid", http.StatusInternalServerError, true).WithOrigin(err)
	}
	return cache, nil
}

func (c *APIKeyCacheClient) Set(keyHash string, value APIKeyCache) *exceptions.Exception {
	client, exception := c.redisClient(keyHash)
	if exception != nil {
		return exception
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return exceptions.New("CacheEncodeFailed", "Cache", "SetAPIKey", "API key cache value could not be encoded", http.StatusInternalServerError, true).WithOrigin(err)
	}
	if err := client.Set(formatAPIKeyCacheKey(keyHash), encoded, c.expiresIn).Err(); err != nil {
		return exceptions.New("CacheWriteFailed", "Cache", "SetAPIKey", "API key cache write failed", http.StatusInternalServerError, true).WithOrigin(err)
	}
	return nil
}
