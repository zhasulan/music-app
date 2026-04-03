# Local Development Guide

## Prerequisites
- Go 1.24+ (build images use Go 1.24 alpine)
- Docker & Docker Compose v2
- Flutter 3.19+ (with Dart 3.3+)
- Make

## Environment variables
- `JWT_SECRET` (shared across services; default `supersecret` for local)
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
  - postgres: 5432
  - redis: 6379
  - minio api: 9000 (console 9001)

## Startup order (backend)
1. `make up` (builds images, runs postgres, redis, minio, services)
2. Run DB migrations (idempotent): `make migrate-all` (creates/updates playlist & library tables; safe to rerun)
3. Seed media objects (once): `make seed-media` (requires `mc` cli). This uploads the mp3s from `deploy/media/tracks/`, removes old wavs, sets public download + CORS for local web playback.
3. Verify health: `curl http://localhost:8080/health`

## Backend troubleshooting
- If Postgres tables missing, confirm `deploy/docker-compose/init.sql` mounted and rerun `make up`.
- Logs: `make logs`
- Tear down: `make down`

## Flutter app
- Location: `mobile/flutter-app`
- Install deps: `make flutter-pub-get`
- Run web (Chrome): `make flutter-run-web` (base URL defaults to http://127.0.0.1:8080)
- Run Android emulator: `make flutter-run` (uses `10.0.2.2:8080`)
- Run iOS simulator/macOS: `flutter run --dart-define API_BASE_URL=http://127.0.0.1:8080`
- Real device: use your host LAN IP, e.g. `--dart-define API_BASE_URL=http://192.168.1.50:8080`

### Auth & API paths (via gateway)
- Register: `POST /api/v1/auth/register {email,password}`
- Login: `POST /api/v1/auth/login {email,password}`
- Refresh: `POST /api/v1/auth/refresh {refresh_token}`
- Profile: `GET/PUT /api/v1/users/me`
- Catalog: `GET /api/v1/catalog/tracks`
- Playlists: `POST/GET/PATCH/DELETE /api/v1/playlists`, add/remove tracks
- Library: `PUT/DELETE /api/v1/library/liked-tracks/:trackId`, list
- Playback: `POST /api/v1/playback/{start|pause|resume|seek}`, `GET /api/v1/playback/current`
- Media: `GET /api/v1/media/tracks/:trackId`, `GET /api/v1/media/tracks/:trackId/source`

### Emulator/Simulator networking
- Android emulator: host mapped to `10.0.2.2`
- iOS simulator: `127.0.0.1`
- Physical devices: use host machine LAN IP and ensure firewall allows inbound 8080/8001-8003
- Flutter Web (Chrome): served from localhost random port; API_BASE_URL should be `http://127.0.0.1:8080`

### Sample accounts / data
- No seed users; register a new account via app or `curl`.
- Catalog uses in-memory seed (3 tracks, 2 artists, 2 albums).

### Browser smoke test
1) `make up`
2) `curl http://127.0.0.1:8080/health` (expect `{ "status": "ok" }`)
3) `make flutter-run-web` (opens Chrome)
4) Register a user, then login
5) Open Profile tab (edit/save name & country)
6) Open Catalog tab (see seeded tracks)
7) Logout via AppBar icon

## Commands quick reference
- `make up` / `make down` / `make logs`
- `make flutter-pub-get`
- `make flutter-run`
- `make flutter-analyze`
