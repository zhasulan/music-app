#!/usr/bin/env sh
set -euo pipefail

COMPOSE=${COMPOSE:-"docker compose -f deploy/docker-compose/docker-compose.yml"}
DB=library_service

echo "Ensuring database $DB exists..."
$COMPOSE exec -T postgres psql -U postgres -d postgres -v ON_ERROR_STOP=1 -c "DO \$\$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = '$DB') THEN CREATE DATABASE $DB; END IF; END \$\$;"

echo "Applying library schema..."
$COMPOSE exec -T postgres psql -U postgres -d $DB -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS liked_tracks (
    user_id BIGINT NOT NULL,
    track_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, track_id)
);
SQL

echo "Library migrations completed."
