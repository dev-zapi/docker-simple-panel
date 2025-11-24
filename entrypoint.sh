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
            # Check if target GID is already in use
            if getent group "$DOCKER_SOCK_GID" > /dev/null 2>&1; then
                CONFLICTING_GROUP=$(getent group "$DOCKER_SOCK_GID" | cut -d: -f1)
                echo "Warning: GID $DOCKER_SOCK_GID is already used by group '$CONFLICTING_GROUP'"
                echo "Removing conflicting group to proceed"
                groupdel "$CONFLICTING_GROUP" || echo "Failed to remove conflicting group, continuing anyway"
            fi
            groupmod -g "$DOCKER_SOCK_GID" docker || echo "Warning: Failed to update docker group GID"
        fi
    else
        # Create docker group with the socket's GID
        echo "Creating docker group with GID $DOCKER_SOCK_GID"
        groupadd -g "$DOCKER_SOCK_GID" docker
    fi
    
    # Add appuser to docker group if not already a member
    if getent passwd appuser > /dev/null 2>&1; then
        if ! id -nG appuser 2>/dev/null | grep -qw docker; then
            echo "Adding appuser to docker group"
            usermod -aG docker appuser
        fi
    else
        echo "Warning: appuser does not exist"
    fi
    
    echo "Docker socket permissions configured successfully"
else
    echo "Warning: Docker socket not found at $DOCKER_SOCKET"
    echo "The application may not be able to connect to Docker daemon"
fi

# Execute the main application as appuser
echo "Starting docker-simple-panel..."
exec gosu appuser "$@"
