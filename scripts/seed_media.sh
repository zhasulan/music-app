#!/usr/bin/env sh
set -euo pipefail

MINIO_ENDPOINT=${MINIO_ENDPOINT:-http://localhost:9000}
MINIO_ACCESS_KEY=${MINIO_ACCESS_KEY:-minioadmin}
MINIO_SECRET_KEY=${MINIO_SECRET_KEY:-minioadmin}
BUCKET=${MINIO_BUCKET:-music}
TRACK_DIR=${TRACK_DIR:-deploy/media/tracks}

if ! command -v mc >/dev/null 2>&1; then
  echo "mc (minio client) is required. Install from https://docs.min.io/docs/minio-client-quickstart-guide.html" >&2
  exit 1
fi

mc alias set local "$MINIO_ENDPOINT" "$MINIO_ACCESS_KEY" "$MINIO_SECRET_KEY"
mc mb --ignore-existing local/$BUCKET
mc mirror --overwrite --remove "$TRACK_DIR" local/$BUCKET/tracks

# Allow public read for demo playback (presigned URLs still work)
mc anonymous set download "local/$BUCKET" >/dev/null 2>&1 || true

# Relaxed CORS for local web playback
TMP_CORS="$(mktemp)"
cat >"$TMP_CORS" <<'EOF'
[
  {
    "AllowedMethods": ["GET"],
    "AllowedOrigins": ["*"],
    "AllowedHeaders": ["*"],
    "ExposeHeaders": ["ETag"],
    "MaxAgeSeconds": 3000
  }
]
EOF
mc cors set "local/$BUCKET" "$TMP_CORS"
rm "$TMP_CORS"

mc ls local/$BUCKET/tracks
