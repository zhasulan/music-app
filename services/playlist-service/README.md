# Playlist Service

Manages user playlists.

## Endpoints (behind gateway)
- `GET /health`
- `POST /api/v1/playlists`
- `GET /api/v1/playlists`
- `GET /api/v1/playlists/:id`
- `PATCH /api/v1/playlists/:id`
- `DELETE /api/v1/playlists/:id`
- `POST /api/v1/playlists/:id/tracks`
- `DELETE /api/v1/playlists/:id/tracks/:trackId`

Auth: Bearer JWT via gateway.

## Env
- `HTTP_PORT` (default 8004)
- `DATABASE_URL` (default postgres://postgres:postgres@postgres:5432/playlist_service?sslmode=disable)
- `JWT_SECRET`
- `CATALOG_SERVICE_URL`
- `LOG_LEVEL`

## Database
Tables created via `deploy/docker-compose/init.sql`:
- `playlists (id serial pk, user_id, name, description, timestamps)`
- `playlist_tracks (id serial pk, playlist_id fk, track_id text, position int, added_at timestamp)`

## Run locally
```bash
make run
```

## Notes
- Track existence validated via catalog-service.
- Positions shift automatically on insert/remove.
