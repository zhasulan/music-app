FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
WORKDIR /app/services/recommendation-service
ENV GOWORK=off GOFLAGS=-mod=mod
RUN go build -o /recommendation-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /recommendation-service /usr/local/bin/recommendation-service
EXPOSE 8012
ENTRYPOINT ["/usr/local/bin/recommendation-service"]
