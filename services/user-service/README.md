# User Service

Maintains user profiles and preferences placeholders.

## Endpoints
- `GET /health`
- `GET /api/v1/users/me` (auth required)
- `PUT /api/v1/users/me` (auth required)
- `GET /api/v1/users/me/preferences` (auth required, placeholder)

## Environment
- `HTTP_PORT` (default 8002)
- `LOG_LEVEL` (default info)
- `DATABASE_URL` (default postgres://postgres:postgres@postgres:5432/user_service?sslmode=disable)
- `JWT_SECRET` (must match auth-service)

## Database
Minimal schema:
```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    email TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL DEFAULT '',
    country TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```
- Profiles are auto-created on first authenticated request using the `user_id` from JWT; email defaults to empty until a sync step is added.

## Run locally
```bash
make run
```

## Docker
Built via root `docker compose` (service name `user-service`).
