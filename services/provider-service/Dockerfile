FROM golang:1.24-alpine AS build
WORKDIR /app
# Copy only what we need to resolve modules
COPY shared ./shared
COPY services/provider-service ./services/provider-service
WORKDIR /app/services/provider-service
ENV GOWORK=off GO111MODULE=on
RUN go mod tidy
RUN go mod download
RUN go build -o /provider-service ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /provider-service /usr/local/bin/provider-service
EXPOSE 8013
ENTRYPOINT ["/usr/local/bin/provider-service"]
