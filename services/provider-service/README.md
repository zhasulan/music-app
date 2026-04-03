# provider-service (Audius proxy)

Purpose: expose Audius content through normalized backend endpoints while keeping local/MinIO flows intact.

## Endpoints
- `GET /health`
- `GET /api/v1/providers/audius/search?q=&limit=&offset=` (tracks)
- `GET /api/v1/providers/audius/tracks/trending?limit=&offset=`
- `GET /api/v1/providers/audius/tracks/:id`
- `GET /api/v1/providers/audius/tracks/:id/stream`
- `GET /api/v1/providers/audius/artists/:handle`
- `GET /api/v1/providers/audius/artists/:handle/tracks?limit=&offset=`
- `GET /api/v1/providers/audius/artists/search?q=&limit=&offset=`
- `GET /api/v1/providers/audius/playlists/search?q=&limit=&offset=`
- `GET /api/v1/providers/audius/playlists/:id`
- `GET /api/v1/providers/audius/playlists/:id/tracks?limit=&offset=`

## Env
- `HTTP_PORT` (default 8013)
- `LOG_LEVEL` (info)
- `AUDIOUS_BASE_URL` (default `https://api.audius.co/v1`)
- `AUDIOUS_APP_NAME` (default `FreedomMusic`)
- `CACHE_TTL_SECONDS` (default 300; in-memory)

## Notes
- Read-only proxy; no secrets/tokens required or exposed.
- Responses are normalized with `provider: "audius"` and `providerTrackId`.
- In-memory cache with TTL; no Redis dependency yet.
- Local provider/MinIO remains unchanged.
