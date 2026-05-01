# Rate Limit Library

## Overview

`app/lib/ratelimit` contains Notezy's internal rate-limit utilities:

- a lightweight in-memory limiter (`WeakRateLimiter`),
- a production-oriented hybrid limiter (`HybridRateLimiter`) combining local token bucket checks and cross-server Redis synchronization,
- and a reusable buffer pool utility used by middleware/interceptors.

## Files

- `weak_rate_limiter.go`: simple leaky-bucket style in-memory limiter.
- `hybrid_rate_limiter.go`: hybrid limiter with local + distributed checks and periodic batch sync.
- `reusable_buffer_pool.go`: `sync.Pool`-backed `bytes.Buffer` reuse helper.

## Public API (Key Types)

```go
type WeakRateLimiter
func NewWeakRateLimiter(requestsPerSecond int) *WeakRateLimiter
func (lb *WeakRateLimiter) Allow() bool

type HybridRateLimiter
func NewHybridRateLimiter(
	rateLimit rate.Limit,
	burst int,
	userLimit int32,
	windowDuration time.Duration,
	backendServerName types.BackendServerName,
	isAuthorizedLimiter bool,
) *HybridRateLimiter
func (hrl *HybridRateLimiter) AllowByFingerprint(fingerprint string) (bool, int32)
func (hrl *HybridRateLimiter) AllowByUserId(userId uuid.UUID) (bool, int32)
func (hrl *HybridRateLimiter) Allow(key string) (bool, int32)
func (hrl *HybridRateLimiter) AllowN(key string, now time.Time, n int) (bool, int32)
func (hrl *HybridRateLimiter) GetStatus() map[string]interface{}
func (hrl *HybridRateLimiter) GetDetailStatus() map[string]interface{}
func (hrl *HybridRateLimiter) Stop()

type ReusableBufferPool
func NewReusableBufferPool() *ReusableBufferPool
func (p *ReusableBufferPool) Get() *bytes.Buffer
func (p *ReusableBufferPool) Put(buffer *bytes.Buffer)
```

## Hybrid Limiter Flow

1. Local token bucket (`x/time/rate`) performs fast in-process throttling.
2. Global usage is checked from Redis-backed records across backend servers.
3. Consumed tokens are batched into pending tasks.
4. A periodic sync loop flushes pending deltas to Redis.
5. On sync failure, pending tasks are retried up to configured retry count.

## Example (Hybrid, Authorized)

```go
limiter := ratelimit.NewHybridRateLimiter(
	rate.Limit(10),
	20,
	100,
	time.Minute,
	types.BackendServerName("api-1"),
	true,
)
defer limiter.Stop()

allowed, remaining := limiter.Allow("7b2b77ff-ec50-4f1f-b309-8f6cf8f0f2c7")
if !allowed {
	// reject request (too many requests)
}
_ = remaining
```

## Example (Reusable Buffer Pool)

```go
pool := ratelimit.NewReusableBufferPool()
buf := pool.Get()
buf.Reset()

buf.WriteString("temporary response")

buf.Reset()
pool.Put(buf)
```

## Project Usage Example

- Authorized middleware: `app/middlewares/authorized_rate_limit_middleware.go`
- Unauthorized middleware: `app/middlewares/unauthorized_rate_limit_middleware.go`
- Buffer pool reused in timeout/interceptor logic.

## File Structure

```text
app/lib/ratelimit/
├── README.md
├── LICENSE.md
├── hybrid_rate_limiter.go
├── reusable_buffer_pool.go
└── weak_rate_limiter.go
```
