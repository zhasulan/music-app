# Freedom Music Backend (Phase 3)

Monorepo for a Spotify-like backend plus Flutter client. Phase 3 adds search, events collection, and simple recommendations on top of the Phase 1–2 stack.

## Current Phase Status
- Implemented and runnable locally via Docker Compose.
- Data: Postgres (auth, user, playlist, library, events), Redis (playback sessions), MinIO (media). Catalog remains in-memory seed. No OpenSearch/Kafka.
- Flutter app: discover, search, catalog, playlists, library, playback; works on web/mobile with just_audio.

## Repository Structure
- `services/` — api-gateway, auth-service, user-service, catalog-service, playlist-service, library-service, playback-service, media-service, events-service, search-service, recommendation-service
- `shared/` — config/logger/errors/httpx/middleware/auth/observability
- `deploy/docker-compose/` — docker-compose.yml, init.sql
- `docs/` — local-development.md, phase-3-testing.md, local-smoke-test.md
- `scripts/` — migrate_* (playlist/library/events), seed_media.sh
- `mobile/flutter-app/` — Flutter client
- `Makefile`, `go.work`

## Implemented Services
- **api-gateway**: Reverse proxy for all services under `/api/v1`, health `/health`, CORS + request logging.
- **auth-service**: Register/login/refresh (JWT HS256), Postgres users.
- **user-service**: Profile get/update, Postgres.
- **catalog-service**: Seeded tracks/albums/artists in memory.
- **playlist-service**: Playlists CRUD + tracks, Postgres.
- **library-service**: Like/unlike/list liked tracks, Postgres.
- **playback-service**: Start/pause/resume/seek/current; Redis per-user session.
- **media-service**: Media metadata + MinIO URLs for mp3 demo files.
- **events-service**: Ingests track/search/playlist events into Postgres (no Kafka).
- **search-service**: Substring search + suggestions over catalog seed (no OpenSearch; reload on restart).
- **recommendation-service**: Rule-based trending/recently-played/for-you built from events + catalog (no ML).
- **provider-service**: Audius proxy provider; normalized endpoints for Audius tracks/artists/playlists/streaming (local MinIO flow unchanged).

## Current Features
- Auth, profile, playlists, library, playback with MinIO mp3s.
- Events capture: track_played/paused/liked, playlist_created, playlist_track_added, search_performed.
- Search with suggestions over seeded catalog.
- Recommendations: trending (most played/liked), recently played (per user), for-you (liked + trending blend).
- Flutter: tabs for Discover, Search, Catalog, Playlists, Library, Profile; bottom mini-player; actions play/like/add-to-playlist.
- Structured JSON logs; request-id; CORS enabled for local dev.

## Tech Stack
- Go 1.24, Gin, PostgreSQL 16, Redis 7, MinIO, Docker Compose
- Flutter 3.x, Riverpod, Dio, just_audio, go_router
- No OpenSearch/Kafka in local setup; Audius read-only via provider-service

## Local Development (summary)
See `docs/local-development.md` for full details.

```bash
make up            # build & start all services + db/redis/minio
make migrate-all   # playlist, library, events schemas
make seed-media    # upload demo mp3s to MinIO
make reindex-search # restart search-service (catalog seed reload)
make logs          # follow logs
make down          # stop stack
```
Health: `curl http://127.0.0.1:8080/health`

## Key Environment Variables
- Shared: `JWT_SECRET`, `POSTGRES_PASSWORD`
- api-gateway: service URLs for auth/user/catalog/playlist/library/playback/media/events/search/recommendation
- provider-service: `AUDIOUS_BASE_URL`, `AUDIOUS_APP_NAME`, `CACHE_TTL_SECONDS`
- events-service: `DATABASE_URL` (events_service), `HTTP_PORT`, `JWT_SECRET`
- search-service: `CATALOG_SERVICE_URL`, `HTTP_PORT`
- recommendation-service: `DATABASE_URL` (events_service), `CATALOG_SERVICE_URL`, `JWT_SECRET`, `HTTP_PORT`
- media-service: MinIO endpoint/keys/bucket/public URL

## Smoke Tests (quick)
```bash
curl http://127.0.0.1:8080/health
curl "http://127.0.0.1:8080/api/v1/catalog/tracks"
curl "http://127.0.0.1:8080/api/v1/search?q=love"
curl http://127.0.0.1:8080/api/v1/recommendations/trending
curl http://127.0.0.1:8080/api/v1/providers/audius/tracks/trending
# audius stream (get url):
curl http://127.0.0.1:8080/api/v1/providers/audius/tracks/<id>/stream
# needs token:
curl -H "Authorization: Bearer <token>" -X POST \
  http://127.0.0.1:8080/api/v1/events/track-played \
  -H 'Content-Type: application/json' -d '{"track_id":"trk_1"}'
```

## Running Flutter
```bash
make flutter-pub-get
make flutter-run-web     # Chrome, API_BASE_URL=http://127.0.0.1:8080
make flutter-run         # Android emulator, API_BASE_URL=http://10.0.2.2:8080
```
See `mobile/flutter-app/README.md` for platform notes.

## Search Notes
- No OpenSearch; search-service loads catalog seed into memory at startup.
- Restart search-service (or `make reindex-search`) after catalog seed changes.

## Recommendation Notes
- Heuristic only: counts of plays/likes and recent user history; no ML/vector search.
- Depends on events being recorded.

## Event Notes
- Stored in Postgres (events_service DB); no Kafka pipeline.
- Endpoints: generic `/api/v1/events` plus typed track-played/paused/liked, playlist-created/track-added, search_performed.

## Makefile Commands
- `make up` / `make down` / `make logs` / `make restart` / `make rebuild`
- `make migrate-all` (playlist, library, events)
- `make seed-media`
- `make reindex-search` (restart search-service)
- `make test` / `make lint`
- Flutter: `make flutter-pub-get`, `make flutter-run`, `make flutter-run-web`, `make flutter-analyze`

## Limitations
- Catalog still in-memory seed (no DB, no images).
- Search is substring/in-memory; no OpenSearch.
- Recommendations are rule-based; no ML.
- No subscriptions/billing; no production hardening (K8s/CDN/observability stack).
- Tokens are not revoked; refresh flow basic.
- No Kafka/event streaming; events only in Postgres.

## Next Planned Phase
- Possible future: richer catalog data store, analytics/observability, ML recommendations, social features, production deployment hardening.
