FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/media-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /media-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /media-service /usr/local/bin/media-service
EXPOSE 8007
ENTRYPOINT ["/usr/local/bin/media-service"]
