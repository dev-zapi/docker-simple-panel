# docker-simple-panel
A simple docker containers dashboard.

## Overview
This is a Go backend application that provides a REST API for managing Docker containers. It includes user authentication, container status monitoring, and container operations.

## Features
- User authentication with JWT tokens
- User registration and login
- SQLite database for user management
- Docker container listing with health status
- Container operations: start, stop, restart
- Docker daemon connectivity via `/var/run/docker.sock`

## Requirements
- Go 1.21 or higher (for building from source)
- Docker daemon running
- SQLite3 (for building from source)

## Quick Start with Docker

The easiest way to run docker-simple-panel is using the pre-built Docker image:

```bash
# Pull the image from GitHub Container Registry
docker pull ghcr.io/dev-zapi/docker-simple-panel:latest

# Run the container
docker run -d \
  --name docker-simple-panel \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./data:/app/data \
  -e JWT_SECRET=CHANGE_ME_STRONG_SECRET \
  ghcr.io/dev-zapi/docker-simple-panel:latest

# Check if the service is running
curl http://localhost:8080/api/health
```

**Important Security Notes**:
- The Docker socket (`/var/run/docker.sock`) must be mounted for the application to manage containers.
- **Always change the JWT_SECRET** to a strong, random value in production.
- Consider disabling registration after creating initial users with `-e DISABLE_REGISTRATION=true`.

## Development with Dev Container

For Windows developers or those who prefer containerized development environments, this project includes VS Code Dev Container configuration:

```bash
# Prerequisites
# 1. Install VS Code and Docker Desktop
# 2. Install "Dev Containers" extension in VS Code
# 3. Open project in VS Code
# 4. Press F1 -> "Dev Containers: Reopen in Container"
```

See [.devcontainer/README.md](.devcontainer/README.md) for detailed instructions.

## Installation

### Building from Source

```bash
# Clone the repository
git clone https://github.com/dev-zapi/docker-simple-panel.git
cd docker-simple-panel

# Install dependencies
go mod download

# Build the application
go build -o docker-simple-panel .
```

### Building Docker Image Locally

```bash
# Build the Docker image
docker build -t docker-simple-panel:local .

# Run the container
docker run -d \
  --name docker-simple-panel \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./data:/app/data \
  -e JWT_SECRET=CHANGE_ME_STRONG_SECRET \
  docker-simple-panel:local
```

## Configuration

The application can be configured using environment variables:

- `SERVER_PORT`: Server port (default: 8080)
- `DATABASE_PATH`: Path to SQLite database file (default: ./docker-panel.db)
- `JWT_SECRET`: Secret key for JWT token signing (default: your-secret-key-change-in-production)
- `DOCKER_SOCKET`: Path to Docker socket (default: /var/run/docker.sock)
- `DISABLE_REGISTRATION`: Disable user registration endpoint (default: false, set to "true", "1", or "yes" to disable)

## Running

```bash
# Run with default settings
./docker-simple-panel

# Run with custom configuration
SERVER_PORT=3000 JWT_SECRET=my-secret-key ./docker-simple-panel

# Run with registration disabled
DISABLE_REGISTRATION=true ./docker-simple-panel

# Run with custom docker socket path
DOCKER_SOCKET=/custom/path/docker.sock ./docker-simple-panel
```

## API Endpoints

### API Documentation

The API is fully documented using OpenAPI v3 specification. You can access the OpenAPI JSON documentation at:

```
GET /api/openapi.json
```

This endpoint returns a complete OpenAPI v3 formatted JSON specification that includes:
- All API endpoints with detailed descriptions
- Request/response schemas
- Authentication requirements
- Example values

You can use tools like [Swagger UI](https://swagger.io/tools/swagger-ui/) or [Swagger Editor](https://editor.swagger.io/) to visualize and interact with the API documentation. Simply load the OpenAPI JSON from `http://localhost:8080/api/openapi.json`.

### Public Endpoints

#### Health Check
```
GET /api/health
```
Returns server health status.

#### Register User
```
POST /api/auth/register
Content-Type: application/json

{
  "username": "admin",
  "password": "password123",
  "nickname": "Administrator"
}
```

#### Login
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin",
    "nickname": "Administrator"
  }
}
```

### Protected Endpoints (Require JWT Token)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <token>
```

#### List Containers
```
GET /api/containers
```

Returns a list of all Docker containers with their status and health information.

Response:
```json
{
  "success": true,
  "data": [
    {
      "id": "abc123def456",
      "name": "my-container",
      "image": "nginx:latest",
      "state": "running",
      "status": "Up 2 hours",
      "health": "healthy",
      "created": 1699876543
    }
  ]
}
```

#### Get Container Details
```
GET /api/containers/{id}
```

#### Start Container
```
POST /api/containers/{id}/start
```

#### Stop Container
```
POST /api/containers/{id}/stop
```

#### Restart Container
```
POST /api/containers/{id}/restart
```

#### Docker Health Check
```
GET /api/docker/health
```

Checks if the Docker daemon is accessible.

#### Get System Configuration
```
GET /api/config
```

Returns the current system configuration (Docker socket path and registration status).

Response:
```json
{
  "success": true,
  "data": {
    "docker_socket": "/var/run/docker.sock",
    "disable_registration": false
  }
}
```

#### Update System Configuration
```
PUT /api/config
Content-Type: application/json

{
  "docker_socket": "/custom/path/docker.sock",
  "disable_registration": true
}
```

Updates system configuration. Both fields are optional. When `docker_socket` is changed, the Docker client automatically restarts with the new socket path. Configuration persists across server restarts.

Response:
```json
{
  "success": true,
  "message": "Configuration updated successfully",
  "data": {
    "docker_socket": "/custom/path/docker.sock",
    "disable_registration": true
  }
}
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    nickname TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Config Table
```sql
CREATE TABLE config (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
```

Stores system configuration that can be updated at runtime via API:
- `docker_socket`: Docker socket path
- `disable_registration`: Whether user registration is disabled

## Security Notes

- Passwords are hashed using bcrypt before storage
- JWT tokens expire after 24 hours
- Change the default `JWT_SECRET` in production
- The application requires access to `/var/run/docker.sock`

## Development

### Running Tests
```bash
go test ./...
```

### Project Structure
```
.
├── main.go              # Application entry point
├── config/              # Configuration management
├── database/            # Database operations
├── docker/              # Docker client wrapper
├── handlers/            # HTTP request handlers
├── middleware/          # HTTP middleware (auth, cors)
└── models/              # Data models
```

## Docker Image Building

### Automated Builds with GitHub Actions

The project includes a GitHub Actions workflow for building and pushing Docker images to GitHub Container Registry.

**To trigger a build:**
1. Go to the "Actions" tab in the GitHub repository
2. Select "Build and Push Docker Image"
3. Click "Run workflow"
4. Optionally specify a custom tag (default: `latest`)
5. Click "Run workflow"

The workflow will:
- Build multi-platform images (linux/amd64, linux/arm64)
- Push to `ghcr.io/dev-zapi/docker-simple-panel`
- Use build caching for faster builds
- Generate image attestations

See [.github/workflows/README.md](.github/workflows/README.md) for more details.

## License
MIT
