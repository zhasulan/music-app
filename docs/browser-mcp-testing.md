# Browser MCP Testing (Phase 1)

## What MCP can test
- HTTP health endpoints via api-gateway
- Auth flows: register, login, refresh (token issuance)
- Protected profile fetch/update
- Catalog list/detail (in-memory seed)
- Gateway routing/CORS
- Optional browser UI via `docs/mcp-smoke.html`

## What MCP cannot fully test
- Native iOS/Android behavior
- Media playback/streaming
- Future services (playlist, library, playback, media, search, recommendations)

## Prerequisites
- Docker + Docker Compose
- Make
- Browser (Chrome) for UI harness

## Start backend
```bash
make up
```
Services: gateway :8080, auth :8001, user :8002, catalog :8003, postgres :5432, redis :6379. Health: `http://127.0.0.1:8080/health`.

## Start browser test target (Flutter Web)
```bash
docker run --rm -p 8081:8081 \
  -v /Users/zhasulan/GolandProjects/freedom-music/mobile/flutter-app:/app \
  -w /app ghcr.io/cirruslabs/flutter:3.19.6 \
  sh -lc "flutter pub get && flutter run -d web-server --target=lib/app/main.dart \
    --web-hostname=0.0.0.0 --web-port=8081 \
    --dart-define API_BASE_URL=http://host.docker.internal:8080"
# open http://localhost:8081
```

## Test URLs
- Health: `http://127.0.0.1:8080/health`
- Auth: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout`
- Profile: `GET/PUT /api/v1/users/me` (Bearer access token required)
- Catalog: `GET /api/v1/catalog/tracks`, `GET /api/v1/catalog/tracks/{id}`, `.../albums`, `.../artists`

## Smoke scenarios (via Flutter Web UI)
1) Open `http://localhost:8081`.
2) Register with test credentials (e.g., `test@example.com`, strong password) via the Auth screen.
3) Login with the same credentials; expect success.
4) Open Profile tab; expect 200 and user data; edit and save name/country.
5) Open Catalog tab; expect seeded track list.
6) Logout; Profile should then be unauthorized (401) when re-entering without token.
7) Optional: attempt protected call after logout to confirm 401.

## Expected results
- Health: 200
- Register: 201/200 with `access_token` + `refresh_token`
- Login: 200 with tokens
- Profile get/update: 200 with `user` object
- Catalog: 200 with `tracks` array
- Unauthorized: 401 when no/invalid token on protected endpoints

## Seeded test user
No pre-seeded users. Use the Register step to create one (e.g., `testuser@example.com / Password123!`).

## Automation hooks (pseudo-steps)
- Navigate to `mcp-smoke.html`
- Click `[data-testid="btn-ping"]` → expect text containing `200`
- Fill `[data-testid="reg-email"]`, `[data-testid="reg-pass"]`, click `[data-testid="btn-register"]` → expect status `200/201`
- Click `[data-testid="btn-profile"]` with tokens set → expect status `200`
- Click `[data-testid="btn-catalog"]` → expect status `200` and list in output
- Click `[data-testid="btn-logout"]` → tokens cleared
- Click `[data-testid="btn-profile"]` again → expect status `401`

## Troubleshooting
- If gateway health fails: `make logs` and check `fm-gateway`, `fm-auth`, `fm-user`, `fm-catalog`, `fm-postgres`.
- If CORS errors: ensure you loaded the page from `http://localhost` (not file://) and gateway is running (CORS is enabled).
- DB init failures: rerun `docker compose -f deploy/docker-compose/docker-compose.yml down -v` then `make up`.
