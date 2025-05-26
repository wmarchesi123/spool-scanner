# Build stage
FROM golang:1.22-alpine AS builder

# Install git and ca-certificates (needed for fetching dependencies)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod file
COPY go.mod ./

# Copy go.sum only if it exists
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o spool-scanner \
    cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/spool-scanner .

# Copy web assets
COPY --from=builder /app/web ./web

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./spool-scanner"]