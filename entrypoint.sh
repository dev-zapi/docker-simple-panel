#!/bin/bash
set -e

# This entrypoint script fixes docker.sock permission issues when running as non-root user
# It dynamically adjusts the docker group GID to match the mounted docker.sock

DOCKER_SOCKET="${DOCKER_SOCKET:-/var/run/docker.sock}"

# Check if docker socket exists and is accessible
if [ -S "$DOCKER_SOCKET" ]; then
    echo "Docker socket found at $DOCKER_SOCKET"
    
    # Get the GID of the docker socket
    DOCKER_SOCK_GID=$(stat -c '%g' "$DOCKER_SOCKET")
    echo "Docker socket GID: $DOCKER_SOCK_GID"
    
    # Check if docker group exists
    if getent group docker > /dev/null 2>&1; then
        CURRENT_DOCKER_GID=$(getent group docker | cut -d: -f3)
        echo "Current docker group GID: $CURRENT_DOCKER_GID"
        
        # If GIDs don't match, update the docker group
        if [ "$CURRENT_DOCKER_GID" != "$DOCKER_SOCK_GID" ]; then
            echo "Updating docker group GID from $CURRENT_DOCKER_GID to $DOCKER_SOCK_GID"
            groupmod -g "$DOCKER_SOCK_GID" docker
        fi
    else
        # Create docker group with the socket's GID
        echo "Creating docker group with GID $DOCKER_SOCK_GID"
        groupadd -g "$DOCKER_SOCK_GID" docker
    fi
    
    # Add appuser to docker group if not already a member
    if ! id -nG appuser | grep -qw docker; then
        echo "Adding appuser to docker group"
        usermod -aG docker appuser
    fi
    
    echo "Docker socket permissions configured successfully"
else
    echo "Warning: Docker socket not found at $DOCKER_SOCKET"
    echo "The application may not be able to connect to Docker daemon"
fi

# Execute the main application as appuser
echo "Starting docker-simple-panel..."
exec gosu appuser "$@"
