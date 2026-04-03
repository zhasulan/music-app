# Auth Service

Handles registration, login, refresh tokens, and JWT issuance.

## Endpoints
- `GET /health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout` (placeholder)

## Environment
- `HTTP_PORT` (default 8001)
- `LOG_LEVEL` (default info)
- `DATABASE_URL` (default postgres://postgres:postgres@postgres:5432/auth_service?sslmode=disable)
- `JWT_SECRET`
- `ACCESS_TOKEN_MINUTES` (default 15)
- `REFRESH_TOKEN_HOURS` (default 720)

## Database
Minimal schema (run inside PostgreSQL):
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Run locally
```bash
make run
```

## Docker
Built via root `docker compose` (service name `auth-service`).
