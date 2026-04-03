# events-service

Receives user/product events and persists them for downstream analytics/recommendations.

## Endpoints
- `GET /health`
- `POST /api/v1/events` (generic, body must include `event_type`)
- `POST /api/v1/events/track-played`
- `POST /api/v1/events/track-paused`
- `POST /api/v1/events/track-liked`
- `POST /api/v1/events/search`
- `POST /api/v1/events/playlist-created`
- `POST /api/v1/events/playlist-track-added`

Auth: Bearer JWT required (same secret as other services).

## Env
- `HTTP_PORT` (default 8010)
- `DATABASE_URL` (default postgres://postgres:postgres@postgres:5432/events_service?sslmode=disable)
- `JWT_SECRET` (default supersecret)
- `LOG_LEVEL` (info)

## Data ownership
- Stores events in Postgres database `events_service`.
- No Kafka/streaming; storage is append-only table with basic indexes.

## Run locally
```bash
make up          # from repo root
make migrate-events
docker compose -f deploy/docker-compose/docker-compose.yml logs events-service
```

## Limitations
- Best-effort ingestion; failures are not retried here.
- No schema evolution/versioning; payload is minimal (user_id, track_id, playlist_id, query, event_type, timestamp).
