# Rate Limits

Devix applies rate and abuse controls globally in `internal/router/router.go`.

## Global request limit

All routes pass through `middleware.RedisRateLimit`.

- Environment variables: `RATE_LIMIT_REQUESTS`, `RATE_LIMIT_WINDOW`
- Defaults: `100` requests per `1m`
- Key when Redis is configured: `rl:ip:<client_ip>`
- Headers on every checked request:
  - `X-RateLimit-Limit`
  - `X-RateLimit-Remaining`
  - `Retry-After` on `429`
- Redis unavailable or not configured: limiter allows requests instead of blocking.

## Auth endpoint limit

Auth routes under `/api/v1/auth/*` also pass through `middleware.AuthRateLimitRedis`.

- Environment variables: `AUTH_RATE_LIMIT_REQUESTS`, `AUTH_RATE_LIMIT_WINDOW`
- Defaults: `5` requests per `1m`
- Key when Redis is configured: `rl:auth:<client_ip>:<path>`
- Applies to:
  - `POST /api/v1/auth/signup`
  - `POST /api/v1/auth/login`
  - `POST /api/v1/auth/refresh`
  - `POST /api/v1/auth/logout`
  - `POST /api/v1/auth/verify-email`
  - `POST /api/v1/auth/forgot-password`
  - `POST /api/v1/auth/reset-password`

## Login protection

Auth routes also use `middleware.LoginProtection`.

- Key when Redis is configured: `login_fail:<client_ip>`
- Threshold: `10` failed auth responses
- Window: `15m`
- Trigger: route returns `401`
- Response after threshold: `429` with `Retry-After`

## Abuse protection

All routes pass through `middleware.AbuseProtection`.

- Threshold: `50` suspicious events
- Window: `10m`
- Key when Redis is configured: `abuse:<client_ip>`
- Suspicious events:
  - Missing `User-Agent`
  - Response status `401`
  - Response status `429`
- Response after threshold: `403`

## Endpoint-specific limits

No endpoint currently registers a separate per-user or per-route category limiter. `middleware.UserRateLimit` exists, but the router does not attach it to any route group yet.

## Operational notes

- Rate limits are Redis-backed sliding windows.
- In local/dev environments without Redis, rate limiting and abuse blocking are effectively disabled.
- Use the global and auth env vars to tune production behavior before adding per-route category limits.
