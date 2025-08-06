# Multi-stage build for Gin API
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app ./cmd/api

FROM alpine:3.20
RUN adduser -D -H appuser && apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /bin/app /app/app
ENV GIN_MODE=release
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/app"]
