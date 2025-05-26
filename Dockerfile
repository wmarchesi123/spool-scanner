# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./
RUN go mod download

COPY . .

# Build for ARM64 explicitly
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build \
    -ldflags="-w -s" \
    -o spool-scanner \
    cmd/server/main.go

# Final stage - debian-slim
FROM debian:bullseye-slim

# Install ca-certificates for HTTPS
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -u 1001 appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/spool-scanner .

# Copy web assets
COPY --from=builder /app/web ./web

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser

EXPOSE 8080

CMD ["./spool-scanner"]