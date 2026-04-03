FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/user-service
ENV GOWORK=off GOFLAGS=-mod=vendor GOPROXY=off GOSUMDB=off
RUN go build -o /user-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /user-service /usr/local/bin/user-service
EXPOSE 8002
ENTRYPOINT ["/usr/local/bin/user-service"]
