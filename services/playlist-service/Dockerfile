FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/playlist-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /playlist-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /playlist-service /usr/local/bin/playlist-service
EXPOSE 8004
ENTRYPOINT ["/usr/local/bin/playlist-service"]
