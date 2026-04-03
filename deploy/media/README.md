# Local Media Assets (Phase 2)

This folder seeds a few demo **mp3** files for the media-service.

- Place audio under `deploy/media/tracks/`.
- Copies are also in the repo root `songs/` folder (renamed to `trk_1.mp3` ...).
- docker-compose mounts this directory into the MinIO container at `/data/music`.
- media-service expects objects:
  - `tracks/trk_1.mp3`
  - `tracks/trk_2.mp3`
  - `tracks/trk_3.mp3`

If you add your own files, keep the same object keys or update the mapping in `services/media-service/internal/service/media_service.go`.

To re-upload after changes:
```
make seed-media
```
