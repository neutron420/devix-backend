# Stage 1: Build
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always 2>/dev/null || echo dev)" \
    -o /devix-backend \
    ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S devix && adduser -S devix -G devix

WORKDIR /app

COPY --from=builder /devix-backend .

RUN chown -R devix:devix /app

USER devix

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./devix-backend"]
