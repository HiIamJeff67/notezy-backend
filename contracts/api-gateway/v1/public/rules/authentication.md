# Authentication and credential rules

APIGateway v1 integrations authenticate with a user-owned `X-API-Key` header. Create the key through the authenticated ClientGateway flow, record the returned secret once, and send it on subsequent API requests.

- The full secret is returned only once. Store it in a secret manager or local environment and never commit it.
- The server persists only a SHA-256 digest and a short display prefix.
- Keys can be expired or revoked; revoked keys fail immediately even if a cache entry exists.
- Do not log request bodies containing credentials, `X-API-Key`, Cookie, Set-Cookie, or CSRF values.
- Unauthorized rate limits are primarily keyed by client IP; API key ID is auxiliary only.
