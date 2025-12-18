# WebUI build stage
FROM node:22-slim AS webui-builder

WORKDIR /app/webui

# Copy webui package files
COPY webui/package.json webui/package-lock.json ./

# Install dependencies
RUN npm ci

# Copy webui source code
COPY webui/ ./

# Create production env file to disable mock API
RUN echo "VITE_USE_MOCK_API=false" > .env.production

# Build webui
RUN npm run build

# Go build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o docker-simple-panel .

# Runtime stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/docker-simple-panel .

# Copy webui from webui-builder
COPY --from=webui-builder /app/webui/dist /app/webui

# Create directory for config
RUN mkdir -p /app/data

# Set environment variables with defaults
ENV CONFIG_PATH=/app/data/config.yaml \
    STATIC_PATH=/app/webui

# Expose port
EXPOSE 8080

# Run the application
CMD ["./docker-simple-panel"]
