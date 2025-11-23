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
- Go 1.21 or higher
- Docker daemon running
- SQLite3

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

```bash
# Clone the repository
git clone https://github.com/dev-zapi/docker-simple-panel.git
cd docker-simple-panel

# Install dependencies
go mod download

# Build the application
go build -o docker-simple-panel .
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

## License
MIT
