# Docker Simple Panel - Copilot Instructions

## Project Overview
Docker Simple Panel is a web application for managing Docker containers with a Go backend and Svelte frontend. It provides user authentication, container management, volume operations, and real-time log streaming.

## Architecture

### Backend (Go)
- **Framework**: Standard library with Gorilla Mux router
- **Database**: SQLite for user management and configuration
- **Docker Integration**: Docker SDK for Go (`github.com/docker/docker`)
- **Authentication**: JWT tokens with 24-hour expiration
- **WebSocket**: Real-time container log streaming

### Frontend (Svelte)
- **Framework**: Svelte 5 with TypeScript
- **Build Tool**: Vite
- **Routing**: svelte-spa-router
- **API**: REST API with JWT authentication

## Project Structure

```
.
├── main.go              # Application entry point, router setup
├── config/              # Configuration management (env vars, database config)
├── database/            # SQLite operations (users, config)
├── docker/              # Docker client wrapper (containers, volumes, manager)
├── handlers/            # HTTP handlers (auth, docker, config, user)
├── middleware/          # HTTP middleware (auth, CORS, logging)
├── models/              # Data models (user, docker containers/volumes)
├── webui/               # Svelte frontend application
│   ├── src/
│   │   ├── components/  # Reusable UI components
│   │   ├── pages/       # Route pages
│   │   ├── services/    # API client services
│   │   ├── stores/      # Svelte stores for state
│   │   └── types/       # TypeScript type definitions
│   └── package.json
└── openapi.json         # API specification
```

## Build and Test Commands

### Backend (Go)
```bash
# Build
go build -o docker-simple-panel .

# Run
./docker-simple-panel

# Run tests (if available)
go test ./...

# With custom config
SERVER_PORT=3000 JWT_SECRET=my-secret ./docker-simple-panel
```

### Frontend (Svelte)
```bash
cd webui

# Development server
npm run dev

# Build for production
npm run build

# Type checking
npm run check
```

### Docker
```bash
# Build multi-stage Docker image (includes frontend build)
docker build -t docker-simple-panel .

# Run container
docker run -d \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./data:/app/data \
  -e JWT_SECRET=CHANGE_ME \
  docker-simple-panel
```

## Code Conventions

### Go Code Style
1. **Error Handling**: Always check and handle errors explicitly
2. **JSON Response**: Use helper functions (`JSONSuccess`, `JSONError`) for consistent API responses
3. **Struct Tags**: Use JSON tags for all exported struct fields
4. **Empty Slices**: Initialize slices as `[]Type{}` instead of `nil` to ensure JSON marshals as `[]` not `null`
5. **Locking**: Docker write operations use `m.mu.Lock()`, read operations use `m.mu.RLock()`
6. **Context**: Use `context.Context` for Docker operations and cancellation

### TypeScript/Svelte Style
1. **TypeScript**: Always use type definitions for API responses
2. **Stores**: Use Svelte stores for shared state (auth, containers)
3. **Components**: Break down UI into reusable components
4. **API Service**: Centralize API calls in `services/api.ts`

### Naming Conventions
- **Handlers**: `<Resource>Handler` struct with `New<Resource>Handler` constructor
- **Models**: Use descriptive names like `ContainerInfo`, `VolumeInfo`
- **API Routes**: RESTful pattern `/api/<resource>/<id>/<action>`
- **Database Methods**: `Get<Resource>`, `Create<Resource>`, `Update<Resource>`

## Security Considerations

1. **Authentication**: All protected routes require JWT token via `Authorization: Bearer <token>` header
2. **Password Hashing**: Use bcrypt for password storage (never store plain text)
3. **JWT Secret**: Must be changed in production via `JWT_SECRET` environment variable
4. **Docker Socket**: Mounting `/var/run/docker.sock` grants full Docker control - use with caution
5. **Input Validation**: Validate all user inputs before processing
6. **CORS**: CORS middleware configured for cross-origin requests
7. **Registration**: Can be disabled via `DISABLE_REGISTRATION=true` environment variable

## Key Patterns

### API Response Format
All API responses follow this structure:
```go
{
  "success": true/false,
  "message": "description",
  "data": <payload>
}
```

### Double Confirmation Pattern
UI uses double confirmation for destructive actions:
- First click sets confirmation state with 3-second timeout
- Second click executes action

### Container Grouping
- Supports three modes: `none`, `compose`, `label`
- Group IDs use prefixes (`label-` for label groups, `_ungrouped_label_` for ungrouped)
- Quick navigation sidebar shown for both compose and label grouping

### WebSocket Log Streaming
- Endpoint: `/api/containers/{id}/logs/stream`
- JWT auth via query parameter or header
- Starts from 30 minutes ago
- Streams both stdout and stderr with timestamps

## Environment Variables

- `SERVER_PORT`: Server port (default: 8080)
- `DATABASE_PATH`: SQLite database path (default: ./docker-panel.db)
- `JWT_SECRET`: JWT signing secret (required in production)
- `DOCKER_SOCKET`: Docker socket path (default: /var/run/docker.sock)
- `DISABLE_REGISTRATION`: Disable user registration (default: false)
- `LOG_LEVEL`: Logging level (error/warn/info/debug, default: info)
- `VOLUME_EXPLORER_IMAGE`: Docker image for volume exploration (default: ghcr.io/dev-zapi/docker-simple-panel:latest)
- `STATIC_PATH`: Path to frontend static files (default: empty for dev, /app/webui in Docker)

## API Endpoints

### Public
- `GET /api/health` - Health check
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login (returns JWT)
- `GET /api/openapi.json` - OpenAPI specification

### Protected (Require JWT)
- `GET /api/containers` - List containers
- `GET /api/containers/{id}` - Container details
- `POST /api/containers/{id}/start` - Start container
- `POST /api/containers/{id}/stop` - Stop container
- `POST /api/containers/{id}/restart` - Restart container
- `GET /api/containers/{id}/logs/stream` - WebSocket log stream
- `GET /api/docker/health` - Docker daemon health
- `GET /api/volumes` - List volumes with container associations
- `DELETE /api/volumes/{name}` - Delete volume
- `GET /api/config` - Get system configuration
- `PUT /api/config` - Update system configuration

## Common Tasks

### Adding a New API Endpoint
1. Define model in `models/` if needed
2. Add handler method to appropriate handler in `handlers/`
3. Register route in `main.go`
4. Add middleware (JWTAuth) if protected
5. Update OpenAPI spec in `openapi.json`

### Adding Frontend Feature
1. Create/update types in `webui/src/types/`
2. Add API call in `webui/src/services/api.ts`
3. Create component in `webui/src/components/` or page in `webui/src/pages/`
4. Update routing if needed in `App.svelte`

### Database Schema Changes
1. Update model in `models/`
2. Add migration logic in `database/database.go`
3. Update relevant handlers

## Dependencies

### Go Modules
- `github.com/docker/docker` - Docker SDK
- `github.com/golang-jwt/jwt/v5` - JWT authentication
- `github.com/gorilla/mux` - HTTP router
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/mattn/go-sqlite3` - SQLite driver
- `golang.org/x/crypto` - Bcrypt password hashing

### NPM Packages
- `svelte` - UI framework
- `svelte-spa-router` - Client-side routing
- `vite` - Build tool
- `typescript` - Type safety

## Testing
- Manual testing via API endpoints
- Docker socket must be accessible for integration testing
- Frontend development uses mock API when `VITE_USE_MOCK_API=true` in `.env`

## Additional Notes
- Volume sorting: Newest first (by CreatedAt descending)
- Container detection: Self-detection via cgroup analysis
- Graceful shutdown: Handles SIGINT and SIGTERM
- HTTP middleware: Custom responseWriter implements http.Hijacker for WebSocket support
- RestartPolicy.Name from Docker API requires string conversion from `container.RestartPolicyMode`
