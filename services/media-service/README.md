# Media Service

Serves media metadata and presigned URLs for track audio via MinIO.

## Endpoints (behind gateway)
- `GET /health`
- `GET /api/v1/media/tracks/:trackId`
- `GET /api/v1/media/tracks/:trackId/source`

Auth: none for Phase 2 (kept simple). Gateway still fronts it.

## Env
- `HTTP_PORT` (default 8007)
- `MINIO_ENDPOINT` (default minio:9000)
- `MINIO_ACCESS_KEY` / `MINIO_SECRET_KEY` (default minioadmin)
- `MINIO_BUCKET` (default music)
- `MINIO_USE_SSL` (default false)
- `MINIO_PUBLIC_URL` (default http://localhost:9000)
- `CATALOG_SERVICE_URL`
- `URL_EXPIRY_MINUTES` (default 120)
- `LOG_LEVEL`

## Media objects
- Expected keys (inside bucket):
  - `tracks/trk_1.wav`
  - `tracks/trk_2.wav`
  - `tracks/trk_3.wav`

## Run locally
```bash
make run
```

## Notes
- Validates track existence via catalog-service.
- Uses presigned GET; for local dev host override is set via `MINIO_PUBLIC_URL`.
- Add files under `deploy/media/tracks/` and run `make seed-media`.
