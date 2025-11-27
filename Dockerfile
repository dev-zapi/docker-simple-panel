# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o docker-simple-panel .

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates libsqlite3-0 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/docker-simple-panel .

# Create directory for database
RUN mkdir -p /app/data

# Set environment variables with defaults
ENV SERVER_PORT=8080 \
    DATABASE_PATH=/app/data/docker-panel.db \
    DOCKER_SOCKET=/var/run/docker.sock \
    DISABLE_REGISTRATION=false \
    STATIC_PATH=

# Expose port
EXPOSE $SERVER_PORT

# Run the application
CMD ["./docker-simple-panel"]
