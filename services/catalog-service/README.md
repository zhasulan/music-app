# Catalog Service

Provides read-only catalog of artists, albums, and tracks (seeded in-memory data for Phase 1).

## Endpoints
- `GET /health`
- `GET /api/v1/catalog/tracks`
- `GET /api/v1/catalog/tracks/:id`
- `GET /api/v1/catalog/albums`
- `GET /api/v1/catalog/albums/:id`
- `GET /api/v1/catalog/artists`
- `GET /api/v1/catalog/artists/:id`

## Environment
- `HTTP_PORT` (default 8003)
- `LOG_LEVEL` (default info)

## Data
Seeded in memory. Swap repository with PostgreSQL later without changing handlers.

## Run locally
```bash
make run
```

## Docker
Built via root `docker compose` (service name `catalog-service`).
