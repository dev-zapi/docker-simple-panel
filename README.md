# docker-simple-panel
A simple docker containers dashboard.

## Overview
This is a Go backend application that provides a REST API for managing Docker containers. It includes authentication with a single user account configured via YAML, container status monitoring, and container operations.

## Features
- Single-user authentication with JWT tokens
- YAML-based configuration management
- Docker container listing with health status
- Container operations: start, stop, restart
- Real-time container log streaming via WebSocket (with 30-minute history)
- Docker volume management with container associations
- Docker daemon connectivity via `/var/run/docker.sock`
- **Progressive Web App (PWA) support** - Install as an app on iPhone, Android, and desktop
- Offline support with service worker caching

## Requirements
- Go 1.21 or higher (for building from source)
- Docker daemon running

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
  ghcr.io/dev-zapi/docker-simple-panel:latest

# Check if the service is running
curl http://localhost:8080/api/health
```

**Default Credentials**:
- Username: `admin`
- Password: `changeme`

**Important**: Change the default password by editing `/app/data/config.yaml` inside the container or mount a custom config file.

**Security Notes**:
- The Docker socket (`/var/run/docker.sock`) must be mounted for the application to manage containers.
- **Always change the default password** in the configuration file.
- **Always change the JWT secret** in production for security.

## Installing as a PWA

Docker Simple Panel supports Progressive Web App (PWA) installation, allowing you to use it as a native app on iPhone, Android, and desktop devices.

### iPhone/iPad Installation

1. Open the application in Safari (http://your-server:8080)
2. Tap the **Share** button (square with arrow pointing up)
3. Scroll down and tap **"Add to Home Screen"**
4. Customize the name if desired and tap **"Add"**
5. The app will now appear on your home screen like a native app

### Android Installation

1. Open the application in Chrome (http://your-server:8080)
2. Tap the menu (three dots) in the top-right corner
3. Tap **"Add to Home screen"** or **"Install app"**
4. Follow the prompts to install
5. The app will appear in your app drawer and home screen

### Desktop Installation

#### Chrome/Edge/Brave
1. Open the application in your browser
2. Look for the install icon in the address bar (plus sign or computer icon)
3. Click the install button
4. The app will open in its own window and appear in your applications

#### Safari (macOS)
1. Open the application in Safari
2. Go to **File** → **Add to Dock**
3. The app will appear in your Dock

### PWA Features
- **Offline Support**: The app caches resources for offline access
- **Fast Loading**: Instant loading from cache
- **App-like Experience**: Runs in standalone mode without browser UI
- **Home Screen Icon**: Custom Docker-themed icon
- **Auto-updates**: Service worker automatically updates the app

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
  docker-simple-panel:local
```

## Configuration

The application uses a YAML configuration file located at `./config.yaml` (or path specified by `CONFIG_PATH` environment variable).

### Configuration File Example

See `config.yaml.example` for a full example. The configuration file includes:

```yaml
# Authentication credentials
username: admin
password: changeme  # Will be automatically hashed with bcrypt

# Server configuration
server:
  port: "8080"
  jwt_secret: "your-secret-key-change-in-production"
  session_max_timeout: 24  # Session timeout in hours

# Docker configuration
docker:
  socket: "/var/run/docker.sock"
  volume_explorer_image: "ghcr.io/dev-zapi/docker-simple-panel:latest"

# Logging configuration
logging:
  level: "info"  # Options: error, warn, info, debug

# Static files (optional)
static_path: ""
```

### Environment Variables

- `CONFIG_PATH`: Path to YAML configuration file (default: ./config.yaml)
- `STATIC_PATH`: Path to static files directory for serving frontend (default: empty for local development; `/app/webui` in Docker image)

### Password Management

- When the application starts for the first time, it creates a default `config.yaml` with username `admin` and password `changeme`
- The password is automatically hashed using bcrypt when saved to the config file
- You can change the password by editing the config file with a plain text password, and it will be hashed on the next startup
- **IMPORTANT**: Always change the default password in production

### Configuration Updates

Most configuration settings can be updated through the web UI Settings page or via the `/api/config` API endpoint. Changes are persisted to the YAML file.

## Running

```bash
# Run with default settings (uses ./config.yaml)
./docker-simple-panel

# Run with custom config path
CONFIG_PATH=/path/to/config.yaml ./docker-simple-panel

# Run with static file serving
STATIC_PATH=/path/to/frontend/dist ./docker-simple-panel
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

#### Login
```
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "changeme"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "username": "admin"
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

#### Stream Container Logs (WebSocket)
```
GET /api/containers/{id}/logs/stream
```

Establishes a WebSocket connection to stream container logs in real-time. The endpoint:
- Requires JWT authentication via the `Authorization: Bearer <token>` header
- Starts streaming from logs generated in the past 30 minutes
- Continues streaming new logs as they are generated
- Includes timestamps for each log line
- Automatically handles both stdout and stderr streams

Example using JavaScript:
```javascript
const token = "your-jwt-token";
const containerId = "abc123def456";
const ws = new WebSocket(
  `ws://localhost:8080/api/containers/${containerId}/logs/stream`,
  [],
  { headers: { Authorization: `Bearer ${token}` } }
);

ws.onmessage = (event) => {
  console.log(event.data); // Each log line with timestamp
};
```

#### Docker Health Check
```
GET /api/docker/health
```

Checks if the Docker daemon is accessible.

#### List Volumes
```
GET /api/volumes
```

Returns a list of all Docker volumes with their associated containers.

Response:
```json
{
  "success": true,
  "data": [
    {
      "name": "my-volume",
      "driver": "local",
      "mountpoint": "/var/lib/docker/volumes/my-volume/_data",
      "created_at": "2025-12-04T02:48:54Z",
      "scope": "local",
      "containers": ["abc123def456"]
    }
  ]
}
```

#### Get System Configuration
```
GET /api/config
```

Returns the current system configuration.

Response:
```json
{
  "success": true,
  "data": {
    "docker_socket": "/var/run/docker.sock",
    "log_level": "info",
    "volume_explorer_image": "ghcr.io/dev-zapi/docker-simple-panel:latest",
    "session_max_timeout": 24,
    "username": "admin"
  }
}
```

#### Update System Configuration
```
PUT /api/config
Content-Type: application/json

{
  "docker_socket": "/custom/path/docker.sock",
  "log_level": "debug",
  "volume_explorer_image": "alpine:latest",
  "session_max_timeout": 48
}
```

Updates system configuration. All fields are optional. When `docker_socket` is changed, the Docker client automatically restarts with the new socket path. Configuration persists to the YAML file.

Response:
```json
{
  "success": true,
  "message": "Configuration updated successfully",
  "data": {
    "docker_socket": "/custom/path/docker.sock",
    "log_level": "debug",
    "volume_explorer_image": "alpine:latest",
    "session_max_timeout": 48,
    "username": "admin"
  }
}
```

## Security Notes

- Passwords are hashed using bcrypt before storage in the config file
- JWT tokens expire after the configured session timeout (default: 24 hours)
- Change the default password and JWT secret in production
- The application requires access to `/var/run/docker.sock`
- Only one user account is supported (configured in config.yaml)

## Development

### Running Tests
```bash
go test ./...
```

### Project Structure
```
.
├── main.go              # Application entry point
├── config/              # Configuration management (YAML-based)
├── docker/              # Docker client wrapper
├── handlers/            # HTTP request handlers
├── middleware/          # HTTP middleware (auth, cors, logging)
├── models/              # Data models
└── webui/               # Frontend application (Svelte)
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
