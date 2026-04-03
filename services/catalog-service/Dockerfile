FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/catalog-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /catalog-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /catalog-service /usr/local/bin/catalog-service
EXPOSE 8003
ENTRYPOINT ["/usr/local/bin/catalog-service"]
