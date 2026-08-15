# Authentication and credential rules

ClientGateway v1 Beta integrations authenticate by calling register or login with an account and password. A successful response sets `accessToken` and `refreshToken` as HttpOnly cookies. Clients must use a cookie jar and return both cookies on subsequent requests.

- Access cookie lifetime: 30 minutes.
- Refresh cookie lifetime: 14 days.
- Production cookies are Secure.
- Access cookie uses SameSite Lax; refresh cookie uses SameSite Strict.
- Access and refresh tokens are never returned in the public JSON body.
- The login/register response exposes a CSRF token. Required operations send it as `X-CSRF-Token`.
- When refresh occurs, read the replacement from the `X-CSRF-Token` response header or `refreshableTokens.newCSRFToken`.
- Do not collect Notezy passwords in an untrusted third-party browser application. The Beta flow is intended for a user's own client or trusted server-side integration.
- Do not log request bodies for login/register, Cookie headers, Set-Cookie headers, or CSRF values.

The server may also accept a Bearer access token, but the public Beta contract is the cookie flow. API-key authentication is not part of v1.
