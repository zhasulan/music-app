# recommendation-service

Rule-based recommendations built from events data.

## Endpoints
- `GET /health`
- `GET /api/v1/recommendations/trending` (no auth)
- `GET /api/v1/recommendations/recently-played` (auth required)
- `GET /api/v1/recommendations/for-you` (auth required)

## Env
- `HTTP_PORT` (default 8012)
- `DATABASE_URL` pointing to events DB (default postgres://postgres:postgres@postgres:5432/events_service?sslmode=disable)
- `CATALOG_SERVICE_URL` (default http://catalog-service:8003)
- `JWT_SECRET`

## Data ownership
- Reads from `events_service` Postgres (events table) and catalog-service for metadata.
- No separate storage.

## Limitations
- Heuristic only (counts of plays/likes + recent user history); no ML or embeddings.
- Requires events to exist; empty if no plays/likes recorded.
