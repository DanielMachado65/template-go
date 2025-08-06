# ---------- Base
FROM golang:1.24-alpine AS base
WORKDIR /app
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
WORKDIR /app
ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"

# instala o Air e garante que esteja acessível
RUN go install github.com/air-verse/air@latest && cp /go/bin/air /usr/local/bin/air

COPY . .
ENV GIN_MODE=debug
EXPOSE 8080

# usa .air.toml se existir; senão roda com defaults
ENTRYPOINT ["sh","-lc","[ -f .air.toml ] && exec air -c .air.toml || exec air"]

# ---------- Builder / Runtime (seu prod atual)
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/api

FROM alpine:3.20 AS runtime
RUN adduser -D -H appuser && apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /bin/app /usr/local/bin/app
ENV GIN_MODE=release
EXPOSE 8080
USER appuser
ENTRYPOINT ["/usr/local/bin/app"]
