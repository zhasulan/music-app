FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/search-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /search-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /search-service /usr/local/bin/search-service
EXPOSE 8011
ENTRYPOINT ["/usr/local/bin/search-service"]
