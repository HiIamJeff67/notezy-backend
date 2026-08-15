# Rate limits and retry rules

ClientGateway emits these headers on rate-limited routes:

- `X-RateLimit-Limit`: request allowance for the active window.
- `X-RateLimit-Remaining`: remaining allowance.
- `X-RateLimit-Reset`: Unix timestamp for the next reset estimate.
- `X-RateLimit-Window`: configured window duration.
- `X-RateLimit-Policy`: currently `hybrid-token-bucket`.

Current ClientGateway v1 routes use an IP/fingerprint limit of 1,000 requests per minute with a 100 requests/second token bucket and burst 10. An authenticated-user limiter exists internally but is not a documented allowance for these routes. These are service limits, not permanent entitlements, and may be lowered during Beta.

On HTTP 429, wait until the reset time and add randomized backoff. Retry only idempotent reads or writes carrying an application-level idempotency guarantee. The current public mutations do not generally expose an idempotency key, so a client must reconcile state before retrying a timed-out mutation.
