# search-service

Simple catalog-backed search with substring matching (no OpenSearch yet).

## Endpoints
- `GET /health`
- `GET /api/v1/search?q=term`
- `GET /api/v1/search/suggest?q=term`

## Env
- `HTTP_PORT` (default 8011)
- `CATALOG_SERVICE_URL` (default http://catalog-service:8003)
- `LOG_LEVEL` (info)

## Notes
- Index is loaded from catalog at startup; restart the service after catalog changes to refresh.
- Ranking is simple: exact/prefix > contains.

## Data ownership
- No dedicated DB; keeps an in-memory index built from catalog-service responses.
- Catalog is the source of truth.

## Limitations
- No OpenSearch/vector/typo tolerance; substring match only.
- Requires service restart to pick up catalog changes.
