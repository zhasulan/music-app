# Playback Service

Stores the active playback session per user in Redis.

## Endpoints (behind gateway)
- `GET /health`
- `POST /api/v1/playback/start` (track_id, position_ms)
- `POST /api/v1/playback/pause` (position_ms)
- `POST /api/v1/playback/resume`
- `POST /api/v1/playback/seek` (position_ms)
- `GET /api/v1/playback/current`

Auth: Bearer JWT via gateway.

## Env
- `HTTP_PORT` (default 8006)
- `REDIS_ADDR` (default redis:6379)
- `REDIS_DB` (default 0)
- `JWT_SECRET`
- `CATALOG_SERVICE_URL`
- `SESSION_TTL_MINUTES` (default 720)
- `LOG_LEVEL`

## Data
- Session stored as JSON in Redis key `playback:{user_id}` with TTL.

## Run locally
```bash
make run
```
