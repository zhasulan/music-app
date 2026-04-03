FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/auth-service
ENV GOWORK=off GOFLAGS=-mod=vendor GOPROXY=off GOSUMDB=off
RUN go build -o /auth-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /auth-service /usr/local/bin/auth-service
EXPOSE 8001
ENTRYPOINT ["/usr/local/bin/auth-service"]
