# Freedom Music Backend (Phase 1)

## Overview
Monorepo for a Spotify-like backend. Phase 1 delivers four Go services (api-gateway, auth-service, user-service, catalog-service), shared libraries, and local Docker Compose. A minimal Flutter client shell lives under `mobile/flutter-app/` for basic auth/profile/catalog flows.

## Current Phase Status
- Phase 2 implemented and runnable locally via Docker Compose.
- Data: Postgres for auth/user; catalog uses in-memory seeded data; Redis available for future caching.
- Flutter: web-capable shell present; runs via Flutter SDK (or container-based command below).

## Repository Structure
- `services/` — api-gateway, auth-service, user-service, catalog-service
- `shared/` — common config/logger/errors/httpx/middleware/auth/observability
- `deploy/docker-compose/` — docker-compose.yml, init.sql
- `docs/` — local-development.md, local-smoke-test.md
- `scripts/` — placeholders
- `mobile/flutter-app/` — Flutter client skeleton (web & mobile targets)
- `Makefile`
- `go.work`

## Implemented Services
- **api-gateway**: Reverse proxy to auth/user/catalog/playlist/library/playback/media under `/api/v1`; health at `/health`; request logging, request-id, CORS (dev-friendly).
- **auth-service**: Register, login, refresh token, logout placeholder; JWT issuance (HS256); Postgres-backed users table.
- **user-service**: Profile fetch/update for current user; auto-creates profile on first access; Postgres-backed `users` table.
- **catalog-service**: Lists tracks/albums/artists and fetch-by-id using in-memory seeded dataset; health at `/health`.
- **playlist-service**: CRUD playlists, add/remove tracks, per-user ownership, Postgres-backed.
- **library-service**: Like/unlike tracks, list liked tracks, Postgres-backed.
- **playback-service**: Stores active playback session per user in Redis; start/pause/resume/seek/current.
- **media-service**: Returns media metadata and presigned URLs from MinIO for tracks; health at `/health`.

## Current Features
- Health endpoints on all services (`/health`).
- Auth: register, login, refresh token; bcrypt password hashing; JWT access/refresh issuance.
- User: get/update current profile (name, country); profiles stored in Postgres.
- Catalog: list tracks/albums/artists and get by id (seeded data, no DB yet).
- Gateway routing for `/api/v1/auth/*`, `/api/v1/users/*`, `/api/v1/catalog/*`.
- Structured JSON logging with request-id; simple CORS for browser local dev.

## Tech Stack
- Go 1.24
- Gin
- PostgreSQL 16
- Redis 7
- Docker Compose
- Flutter 3.x (client shell)
- JWT (HS256), bcrypt, zap logger

## Local Development
See `docs/local-development.md` for details. Quick start:

### Environment Variables
A `.env.example` is provided with base secrets. Key vars (envDefault shown):
- **api-gateway**: `APP_NAME=api-gateway`, `HTTP_PORT=8080`, `LOG_LEVEL=info`, `AUTH_SERVICE_URL`, `USER_SERVICE_URL`, `CATALOG_SERVICE_URL`
- **auth-service**: `APP_NAME=auth-service`, `HTTP_PORT=8001`, `LOG_LEVEL=info`, `DATABASE_URL=postgres://postgres:postgres@postgres:5432/auth_service?sslmode=disable`, `JWT_SECRET`, `ACCESS_TOKEN_MINUTES=15`, `REFRESH_TOKEN_HOURS=720`
- **user-service**: `APP_NAME=user-service`, `HTTP_PORT=8002`, `LOG_LEVEL=info`, `DATABASE_URL=postgres://postgres:postgres@postgres:5432/user_service?sslmode=disable`, `JWT_SECRET`
- **catalog-service**: `APP_NAME=catalog-service`, `HTTP_PORT=8003`, `LOG_LEVEL=info`
- Postgres/Redis handled by docker-compose defaults (user/pass `postgres/postgres`).

### Running the Backend
```bash
make up            # build & start all services + db/redis
make logs          # follow logs
make restart       # restart containers
make down          # stop stack
make rebuild       # rebuild images without cache
```
Health check (host): `curl http://127.0.0.1:8080/health`

### Running Flutter locally
Requirements: Flutter 3.x. If local SDK is broken, use the container command below.
- Android emulator (host Flutter): `make flutter-run` (API_BASE_URL defaults to http://10.0.2.2:8080)
- Web via Chrome (host Flutter): `make flutter-run-web` (API_BASE_URL=http://127.0.0.1:8080)
- Containerized Flutter Web (works without host SDK):
```bash
docker run --rm -p 8081:8081 \
  -v /Users/zhasulan/GolandProjects/freedom-music/mobile/flutter-app:/app \
  -w /app ghcr.io/cirruslabs/flutter:3.19.6 \
  sh -lc "flutter pub get && flutter run -d web-server --target=lib/app/main.dart \
    --web-hostname=0.0.0.0 --web-port=8081 \
    --dart-define API_BASE_URL=http://host.docker.internal:8080"
# then open http://localhost:8081
```

## Smoke Test Scenarios
1) Start stack: `make up`
2) Check gateway health: `curl http://127.0.0.1:8080/health`
3) Register: `curl -X POST http://127.0.0.1:8080/api/v1/auth/register -d '{"email":"a@b.com","password":"password123"}' -H 'Content-Type: application/json'`
4) Login: `curl -X POST http://127.0.0.1:8080/api/v1/auth/login ...` (store access/refresh)
5) Profile: `curl -H "Authorization: Bearer <access>" http://127.0.0.1:8080/api/v1/users/me`
6) Update profile: PUT `/api/v1/users/me` with name/country
7) Catalog list: `curl http://127.0.0.1:8080/api/v1/catalog/tracks`
8) (Optional) Flutter Web: open http://localhost:8081, register/login, view Profile and Catalog tabs, logout.

## Makefile Commands
- `make up` / `make down` / `make logs` / `make restart` / `make rebuild`
- `make build` (docker compose build)
- `make lint` (gofmt current tree)
- `make test` (go test ./...)
- Flutter helpers: `make flutter-pub-get`, `make flutter-run`, `make flutter-run-web`, `make flutter-analyze`

## What is not implemented yet
- Search, recommendations, subscriptions
- Production deployments (Kubernetes, CDNs, etc.)
- Token revocation/session store, advanced caching
- Catalog persistence (currently in-memory seed only)

## Next planned phase
Phase 3 roadmap: events-service, search-service (not implemented).
