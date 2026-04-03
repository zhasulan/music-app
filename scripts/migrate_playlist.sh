#!/usr/bin/env sh
set -euo pipefail

COMPOSE=${COMPOSE:-"docker compose -f deploy/docker-compose/docker-compose.yml"}
DB=playlist_service

echo "Ensuring database $DB exists..."
$COMPOSE exec -T postgres psql -U postgres -d postgres -v ON_ERROR_STOP=1 -c "DO \$\$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = '$DB') THEN CREATE DATABASE $DB; END IF; END \$\$;"

echo "Applying playlist schema..."
$COMPOSE exec -T postgres psql -U postgres -d $DB -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS playlists (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS playlist_tracks (
    id SERIAL PRIMARY KEY,
    playlist_id INT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL,
    position INT NOT NULL,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_playlist_tracks_order ON playlist_tracks (playlist_id, position);
SQL

echo "Playlist migrations completed."
