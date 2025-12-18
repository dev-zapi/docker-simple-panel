# GitHub Actions Workflows

## Docker Build and Push

### Overview
This workflow builds a Docker image and pushes it to GitHub Container Registry (ghcr.io).

### Triggers

The workflow is triggered in two ways:

1. **Automatically**: When a pull request is merged to the `main` branch
   - The image is automatically built and tagged as `latest`
   - Also includes a SHA-based tag with branch prefix

2. **Manually**: Via `workflow_dispatch` in the GitHub Actions UI
   - Allows specifying a custom tag
   - Default tag is `latest` if not specified

### Manual Usage

1. Go to the "Actions" tab in the GitHub repository
2. Select "Build and Push Docker Image" from the workflows list
3. Click "Run workflow"
4. Optionally specify a custom tag (default is `latest`)
5. Click "Run workflow" to start the build

### Parameters
- **tag**: Docker image tag (e.g., `v1.0.0`, `latest`, `dev`). Default: `latest` (only for manual triggers)

### Features
- Multi-platform builds (linux/amd64, linux/arm64)
- Automatic push to GitHub Container Registry
- Build caching using GitHub Actions cache
- Image attestation for supply chain security
- Multi-stage build for minimal image size

### Image Location
After the workflow completes, the image will be available at:
```
ghcr.io/dev-zapi/docker-simple-panel:TAG
```

### Pull the Image
```bash
docker pull ghcr.io/dev-zapi/docker-simple-panel:latest
```

### Run the Container
```bash
docker run -d \
  -p 8080:8080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./data:/app/data \
  -e JWT_SECRET=your-secret-key \
  ghcr.io/dev-zapi/docker-simple-panel:latest
```
