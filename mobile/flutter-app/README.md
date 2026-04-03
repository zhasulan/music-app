# Freedom Music Flutter App (Phase 3)

Playable MVP client for the Freedom Music backend. All traffic goes through the API Gateway.

## Stack
- Flutter 3.19+
- Riverpod
- Dio
- flutter_secure_storage
- go_router
- just_audio
- google_fonts

## Structure
```
lib/
  app/            # main app + routing
  core/           # config, API client, token storage, models
  features/
    auth/         # login/register/splash
    profile/      # profile view/edit
    catalog/      # catalog list + play/add/like
    playback/     # mini player + controller (just_audio)
    library/      # liked tracks
    playlists/    # playlists list/detail/add track
    media/        # media source client
```

## Run
```bash
# install deps
flutter pub get
# web (Chrome)
flutter run -d chrome --dart-define API_BASE_URL=http://127.0.0.1:8080
# Android emulator (host at 10.0.2.2)
flutter run --dart-define API_BASE_URL=http://10.0.2.2:8080
# iOS simulator / macOS
flutter run --dart-define API_BASE_URL=http://127.0.0.1:8080
# real device (replace with your LAN IP)
flutter run --dart-define API_BASE_URL=http://192.168.1.50:8080
```

## Screens & flows
- Splash (token check)
- Auth: Login / Register (polished forms)
- Home tabs: Discover, Search, Catalog, Playlists, Library, Profile
- Audius tab: Audius trending (external provider, read-only)
- Playback: bottom mini-player (play/pause/seek)
- Playlists: create/list/detail, add/remove tracks
- Library: like/unlike tracks, list liked
- Profile: view/update basic fields
- Discover: trending / recently played / for you (recommendation-service)
- Search: tracks/albums/artists with suggestions and play/add/like actions
- Events: app sends play/pause/like/playlist/search events to events-service (best-effort)
- Logout (AppBar icon)
- Uses real mp3 demo files from MinIO via media-service (run `make seed-media` after `make up`)
- Audius: uses provider-service proxy; shows provider badge and streams via Audius URLs (no secrets in app)

## Notes
- All API calls go through the gateway.
- Playback uses media-service for signed URLs and playback-service for state.
- Tokens stored in secure storage; refresh flow still placeholder.
- Web works; audio on web relies on just_audio_web (served via same API base).
- For local stacks, run `make up` then `make seed-media` to ensure media objects exist.

## Known limitations
- No advanced queue/background audio.
- Refresh token flow remains minimal.
- Catalog metadata is basic (no images yet).
- Recommendations are rule-based; results depend on recorded events.
- Search is substring-based; suggestions are simple and limited.
- Event sends are fire-and-forget; failures are not surfaced to the user.
