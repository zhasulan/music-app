APP=catalog-service

run:
	go run ./cmd/server

build:
	go build -o bin/$(APP) ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy
