FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/library-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /library-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /library-service /usr/local/bin/library-service
EXPOSE 8005
ENTRYPOINT ["/usr/local/bin/library-service"]
