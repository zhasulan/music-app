# API Gateway

Entry point for the Freedom Music backend. Routes traffic to downstream services.

## Endpoints
- `GET /health`
- `POST /api/v1/auth/*` → auth-service
- `GET/PUT /api/v1/users/*` → user-service
- `GET /api/v1/catalog/*` → catalog-service

## Environment
- `HTTP_PORT` (default 8080)
- `LOG_LEVEL` (default info)
- `AUTH_SERVICE_URL` (default http://auth-service:8001)
- `USER_SERVICE_URL` (default http://user-service:8002)
- `CATALOG_SERVICE_URL` (default http://catalog-service:8003)

## Run locally
```bash
make run
```

## Docker
Built via root `docker compose` (service name `api-gateway`).
