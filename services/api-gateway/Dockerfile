FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/api-gateway
ENV GOWORK=off GOFLAGS=-mod=vendor GOPROXY=off GOSUMDB=off
RUN go build -o /api-gateway ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /api-gateway /usr/local/bin/api-gateway
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/api-gateway"]
