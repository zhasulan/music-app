FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/playback-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /playback-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /playback-service /usr/local/bin/playback-service
EXPOSE 8006
ENTRYPOINT ["/usr/local/bin/playback-service"]
