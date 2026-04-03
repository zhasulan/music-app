# Local Development Guide (Phase 3)

## Prerequisites
- Go 1.24+
- Docker & Docker Compose v2
- Flutter 3.19+ (Dart 3.3+)
- Make

## Environment variables
- `JWT_SECRET` (default `supersecret` for local)
- `POSTGRES_PASSWORD` (default `postgres`)
- Service ports:
  - api-gateway: 8080
  - auth-service: 8001
  - user-service: 8002
  - catalog-service: 8003
  - playlist-service: 8004
  - library-service: 8005
  - playback-service: 8006
  - media-service: 8007
  - events-service: 8010
  - search-service: 8011
- recommendation-service: 8012
- provider-service: 8013
  - postgres: 5432
  - redis: 6379
  - minio api: 9000 (console 9001)

## Startup order (backend)
1. `make up` (builds images, runs postgres, redis, minio, all services)
2. `make migrate-all` (idempotent: playlist, library, events)
3. `make seed-media` (uploads mp3s from `deploy/media/tracks/`, sets public/CORS)
4. Verify health: `curl http://127.0.0.1:8080/health`

## Verifications
- Search: `curl "http://127.0.0.1:8080/api/v1/search?q=love"` (expects results list; substring-based)
- Suggestions: `curl "http://127.0.0.1:8080/api/v1/search/suggest?q=l"`
- Events: `curl -H "Authorization: Bearer <token>" -X POST http://127.0.0.1:8080/api/v1/events/track-played -d '{"track_id":"trk_1"}' -H 'Content-Type: application/json'`
- Recommendations (no auth for trending): `curl http://127.0.0.1:8080/api/v1/recommendations/trending`
- Recommendations (auth): `curl -H "Authorization: Bearer <token>" http://127.0.0.1:8080/api/v1/recommendations/recently-played`

## Reindex / refresh search
- Search loads catalog seed at service startup. If seed changes, restart: `make reindex-search` (restars search-service).

## Backend troubleshooting
- Missing tables: ensure `deploy/docker-compose/init.sql` mounts; rerun `make up` then `make migrate-all`.
- Media 403/404: rerun `make seed-media`; ensure MinIO is up (9000/9001).
- Gateway 404: check gateway env URLs and container health (`docker compose -f deploy/docker-compose/docker-compose.yml ps`).
- Events DB: if migrations not run, `make migrate-events` (needs docker socket access).

## Flutter app
- Location: `mobile/flutter-app`
- Install deps: `make flutter-pub-get`
- Run web (Chrome): `make flutter-run-web` (API_BASE_URL=http://127.0.0.1:8080)
- Run Android emulator: `make flutter-run` (API_BASE_URL=http://10.0.2.2:8080)
- Run iOS simulator/macOS: `flutter run --dart-define API_BASE_URL=http://127.0.0.1:8080`
- Real device: use host LAN IP, e.g. `--dart-define API_BASE_URL=http://192.168.1.50:8080`

### Auth & API paths (via gateway)
- Register: `POST /api/v1/auth/register {email,password}`
- Login: `POST /api/v1/auth/login {email,password}`
- Profile: `GET/PUT /api/v1/users/me`
- Catalog: `GET /api/v1/catalog/tracks`
- Playlists: `POST/GET/PATCH/DELETE /api/v1/playlists`, add/remove tracks
- Library: `PUT/DELETE /api/v1/library/liked-tracks/:trackId`, list
- Playback: `POST /api/v1/playback/{start|pause|resume|seek}`, `GET /api/v1/playback/current`
- Media: `GET /api/v1/media/tracks/:trackId`, `GET /api/v1/media/tracks/:trackId/source`
- Events: `POST /api/v1/events/*` (track-played/paused/liked, playlist-created/track-added, search)
- Search: `GET /api/v1/search`, `GET /api/v1/search/suggest`
- Recommendations: `GET /api/v1/recommendations/trending`, `recently-played`, `for-you`

### Emulator/Simulator networking
- Android emulator: host is `10.0.2.2`
- iOS simulator: `127.0.0.1`
- Physical devices: use host LAN IP; open firewall for 8080/8001-8012/9000
- Flutter Web: API_BASE_URL should be `http://127.0.0.1:8080`

### Sample data
- No seed users; register manually.
- Catalog seed: 3 tracks, 2 artists, 2 albums in-memory.

### Browser smoke test
1) `make up`
2) `curl http://127.0.0.1:8080/health`
3) `make flutter-run-web`
4) Register + login
5) Discover tab: trending appears; Search tab: query “love” shows results.
6) Play a track; mini-player should show and play audio.
7) Like a track; Library tab shows it.

## Commands quick reference
- `make up` / `make down` / `make logs` / `make restart` / `make rebuild`
- `make migrate-all`
- `make seed-media`
- `make reindex-search`
- `make up provider-service` if only Audius proxy changes
- `make flutter-pub-get` / `make flutter-run` / `make flutter-run-web` / `make flutter-analyze`
- Audius proxy (via provider-service):
  - Trending tracks: `curl http://127.0.0.1:8080/api/v1/providers/audius/tracks/trending`
  - Track stream URL: `curl http://127.0.0.1:8080/api/v1/providers/audius/tracks/<id>/stream`
  - Artist: `curl http://127.0.0.1:8080/api/v1/providers/audius/artists/<handle>`
  - Playlist search: `curl "http://127.0.0.1:8080/api/v1/providers/audius/playlists/search?q=mix"`
