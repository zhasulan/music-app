# AGENTS.md

## Project Overview

This repository is a monorepo for a Spotify-like music streaming backend.

The system should start as a practical MVP-friendly backend, but the structure must be ready for future scaling into a more mature microservice architecture.

This is not a toy project. The generated code should look production-oriented:
- clean project structure
- clear service boundaries
- health checks
- graceful shutdown
- structured logging
- configuration via environment variables
- Docker support
- consistent internal conventions

The repository is a monorepo and all backend services must be placed under one common folder:

- `services/`

Every service lives inside:
- `services/<service-name>`

Shared code must live inside:
- `shared/`

Infrastructure and local environment assets must live in dedicated folders:
- `deploy/`
- `scripts/`
- `docs/`

---

## Main Product Goal

Build a backend platform for a music streaming application similar to Spotify.

Core features for MVP:
- authentication
- users
- music catalog
- playlists
- liked songs / library
- playback session management
- media access for streaming
- event collection for analytics
- search
- foundation for recommendations later

Non-goals for the first stage:
- full ML recommendation engine
- advanced social features
- offline encrypted downloads
- DRM-heavy implementation
- real-time collaborative playlists
- podcast platform
- full multi-region deployment

---

## Architecture Principles

1. Start practical, not overengineered
2. Prefer simple service boundaries with clean code
3. Design for future extraction and scaling
4. Keep each service independently runnable
5. Prefer explicit code over magic abstractions
6. Keep APIs and contracts easy to understand
7. Avoid deep coupling between services
8. Use async events only where it actually helps
9. Build MVP first, then evolve

---

## Monorepo Structure Rules

The repository should follow this shape:

.
├── AGENTS.md
├── README.md
├── docs/
├── deploy/
│   ├── docker-compose/
│   └── k8s/
├── scripts/
├── shared/
│   ├── config/
│   ├── httpx/
│   ├── logger/
│   ├── errors/
│   ├── middleware/
│   ├── auth/
│   ├── events/
│   └── observability/
├── services/
│   ├── api-gateway/
│   ├── auth-service/
│   ├── user-service/
│   ├── catalog-service/
│   ├── playlist-service/
│   ├── library-service/
│   ├── playback-service/
│   ├── media-service/
│   ├── events-service/
│   ├── search-service/
│   └── recommendation-service/
└── Makefile

Important:
- all backend services must stay under `services/`
- shared reusable code goes only into `shared/`

---

## Preferred Tech Stack

- Go 1.24+
- Gin
- PostgreSQL
- Redis
- Kafka
- MinIO
- OpenSearch
- Docker Compose
- OpenTelemetry

---

## Service List

api-gateway:
- entry point
- auth middleware
- routing

auth-service:
- registration
- login
- JWT
- sessions

user-service:
- profile
- preferences
- region

catalog-service:
- tracks
- artists
- albums

playlist-service:
- playlists CRUD
- tracks management

library-service:
- liked tracks
- saved albums

playback-service:
- playback state
- current track

media-service:
- signed URLs
- audio access

events-service:
- event ingestion
- Kafka publishing

search-service:
- search
- autocomplete

recommendation-service:
- trending
- basic recommendations

---

## Development Strategy

Phase 1:
- api-gateway
- auth-service
- user-service
- catalog-service

Phase 2:
- playlist-service
- library-service
- playback-service
- media-service

Phase 3:
- events-service
- search-service

Phase 4:
- recommendation-service

---

## Internal Service Structure

services/<service-name>/
├── cmd/server/main.go
├── internal/
│   ├── domain/
│   ├── repository/
│   ├── service/
│   ├── handler/
│   └── app/
├── Dockerfile
├── Makefile
└── go.mod

---

## Shared Package Rules

Shared code goes into `shared/`

Allowed:
- config
- logger
- errors
- middleware
- auth utils
- event contracts

Do not put business logic here

---

## Data Ownership

Each service owns its data

---

## API Rules

- REST
- JSON
- versioned endpoints `/api/v1`
- health endpoint `/health`

---

## Error Format

{
"error": {
"code": "invalid_request",
"message": "error message"
}
}

---

## Logging

- structured logs
- include request id
- include status code

---

## Configuration

- environment variables only

---

## Docker

- each service must have Dockerfile
- root docker-compose

---

## First Generation Goal

Generate:
- monorepo structure
- shared packages
- api-gateway
- auth-service
- user-service
- catalog-service
- docker-compose
- Makefile

---

## Important

Do not overengineer  
Keep it practical  
Keep it clean  


## Mobile Client

A Flutter mobile application may be added later.

Location:
- mobile/flutter-app/

Responsibilities:
- consume backend APIs via API Gateway
- handle authentication
- playback audio using direct media URLs
- manage UI and user interaction

Important:
- Flutter app must NOT stream audio through backend
- audio must be played directly from object storage/CDN