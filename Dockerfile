# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies (with GOPROXY fallback)
RUN GOPROXY=https://goproxy.io,direct go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o docker-simple-panel .

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies including gosu for privilege dropping
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 gosu && \
    rm -rf /var/lib/apt/lists/*

# Create a non-root user
RUN useradd -m -u 1000 appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/docker-simple-panel .

# Copy entrypoint script
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Create directory for database and set permissions
RUN mkdir -p /app/data && chown -R appuser:appuser /app

# Expose port
EXPOSE 8080

# Set environment variables with defaults
ENV SERVER_PORT=8080 \
    DATABASE_PATH=/app/data/docker-panel.db \
    DOCKER_SOCKET=/var/run/docker.sock \
    DISABLE_REGISTRATION=false

# Set entrypoint to handle docker socket permissions
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

# Run the application
CMD ["./docker-simple-panel"]
