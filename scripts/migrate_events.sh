#!/usr/bin/env sh
set -euo pipefail

COMPOSE=${COMPOSE:-"docker compose -f deploy/docker-compose/docker-compose.yml"}
DB=events_service

echo "Ensuring database $DB exists..."
$COMPOSE exec -T postgres psql -U postgres -d postgres -v ON_ERROR_STOP=1 -c "DO \$\$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = ) THEN CREATE DATABASE $DB; END IF; END \$\$;"

echo "Applying events schema..."
$COMPOSE exec -T postgres psql -U postgres -d $DB -v ON_ERROR_STOP=1 <<SQL
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    event_type TEXT NOT NULL,
    track_id TEXT,
    playlist_id INT,
    query TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events (created_at);
CREATE INDEX IF NOT EXISTS idx_events_user ON events (user_id);
CREATE INDEX IF NOT EXISTS idx_events_type ON events (event_type);
SQL

echo "Events migrations completed."
