COMPOSE=docker compose -f deploy/docker-compose/docker-compose.yml

up:
	$(COMPOSE) up -d --build

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

build:
	$(COMPOSE) build

rebuild:
	$(COMPOSE) build --no-cache

restart:
	$(COMPOSE) restart

test:
	go test ./...

lint:
	gofmt -w $$(find services shared -name '*.go')

seed-media:
	MINIO_ENDPOINT=http://localhost:9000 MINIO_BUCKET=music ./scripts/seed_media.sh

migrate-playlist:
	./scripts/migrate_playlist.sh

migrate-library:
	./scripts/migrate_library.sh

migrate-all: migrate-playlist migrate-library

flutter-pub-get:
	cd mobile/flutter-app && flutter pub get

flutter-run:
	cd mobile/flutter-app && flutter run --dart-define API_BASE_URL=http://10.0.2.2:8080

flutter-analyze:
	cd mobile/flutter-app && flutter analyze

flutter-run-web:
	cd mobile/flutter-app && flutter run -d chrome --dart-define API_BASE_URL=http://127.0.0.1:8080
