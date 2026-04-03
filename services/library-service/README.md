# Library Service

Handles liked tracks per user.

## Endpoints (behind gateway)
- `GET /health`
- `PUT /api/v1/library/liked-tracks/:trackId`
- `DELETE /api/v1/library/liked-tracks/:trackId`
- `GET /api/v1/library/liked-tracks`

Auth: Bearer JWT via gateway.

## Env
- `HTTP_PORT` (default 8005)
- `DATABASE_URL` (default postgres://postgres:postgres@postgres:5432/library_service?sslmode=disable)
- `JWT_SECRET`
- `CATALOG_SERVICE_URL`
- `LOG_LEVEL`

## Database
`liked_tracks (user_id bigint, track_id text, created_at timestamp, pk(user_id, track_id))`

## Run locally
```bash
make run
```

## Notes
- Track existence validated via catalog-service.
