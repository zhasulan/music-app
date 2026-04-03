# Phase 3 Testing Guide

Practical manual checks for search, discovery, events, and recommendations.

## Preconditions
- Stack running: `make up`
- Migrations: `make migrate-all`
- Media seeded: `make seed-media`
- Have a user token (register/login via app or curl).

## Search
1. Open app Search tab (or curl below).
2. Type “love” (or any substring from seed) — results should show tracks; suggestions chips should appear.
3. Tap a track -> starts playback; mini-player appears.
4. Add to playlist and like from results.
5. Curl check:
   ```bash
   curl "http://127.0.0.1:8080/api/v1/search?q=love"
   curl "http://127.0.0.1:8080/api/v1/search/suggest?q=l"
   ```

## Discover / Home
1. Open Discover tab.
2. Sections:
   - Trending: populated once events exist (plays/likes).
   - Recently played: shows tracks you played (requires auth).
   - For you: blend of liked + trending.
3. Tap play/like/add-to-playlist from cards.

## Trending flow
1. Play a track a few times.
2. Call `curl http://127.0.0.1:8080/api/v1/recommendations/trending` — played track should appear.

## Recently played flow
1. With token:
   ```bash
   curl -H "Authorization: Bearer <token>" http://127.0.0.1:8080/api/v1/recommendations/recently-played
   ```
2. Should list tracks you played most recently.

## For-you flow
1. Like a couple of tracks.
2. Call:
   ```bash
   curl -H "Authorization: Bearer <token>" http://127.0.0.1:8080/api/v1/recommendations/for-you
   ```
3. Should include liked/trending items (rule-based, not ML).

## Event flow
1. Play/pause/like/add-to-playlist/search in the app (best-effort sends events).
2. Verify DB quickly (optional):
   ```bash
   docker compose -f deploy/docker-compose/docker-compose.yml exec -T postgres \
     psql -U postgres -d events_service -c "select event_type, track_id, query, created_at from events order by created_at desc limit 10;"
   ```

## Web smoke
1. `make flutter-run-web`
2. Register/login.
3. Discover + Search tabs load.
4. Play a track; audio should play (served from MinIO public URL).

## Mobile smoke (emulator)
1. `make flutter-run` (API_BASE_URL=http://10.0.2.2:8080).
2. Repeat search, discover, playback.

## Expected results
- Search returns seeded catalog items; suggestions limited and simple.
- Discover sections show data once events exist; trending may be empty until plays/likes recorded.
- Playback works from search/discover/catalog.
- Event posting failures do not block UI.

## Known limitations
- Catalog is in-memory; changing seed requires restarting services.
- Search is substring/in-memory; no OpenSearch ranking.
- Recommendations are heuristic; quality depends on recorded events.
- No Kafka or analytics pipelines; events only in Postgres.
- No offline/background audio; refresh token flow minimal.
- Audius provider is read-only; if Audius API is down, Audius tab may show an error but local content still works.
