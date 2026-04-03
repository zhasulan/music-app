FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/events-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /events-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /events-service /usr/local/bin/events-service
EXPOSE 8010
ENTRYPOINT ["/usr/local/bin/events-service"]
