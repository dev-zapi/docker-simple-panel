#!/bin/bash

# Example usage of docker-simple-panel API
# This script demonstrates how to use the API to manage Docker containers

BASE_URL="http://localhost:8080"

echo "=== Docker Simple Panel API Demo ==="
echo ""

# Step 1: Register a user
echo "Step 1: Register a new user"
curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "secure_password_123",
    "nickname": "System Administrator"
  }' | jq .
echo ""

# Step 2: Login and get token
echo "Step 2: Login and obtain JWT token"
TOKEN=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "secure_password_123"
  }' | jq -r '.data.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    echo "✓ Successfully logged in"
    echo "Token (first 50 chars): ${TOKEN:0:50}..."
else
    echo "✗ Login failed"
    exit 1
fi
echo ""

# Step 3: Check Docker daemon health
echo "Step 3: Check Docker daemon connectivity"
curl -s "$BASE_URL/api/docker/health" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Step 4: List all containers
echo "Step 4: List all Docker containers"
curl -s "$BASE_URL/api/containers" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Example: If you have containers running, you can control them
# Uncomment and replace CONTAINER_ID with actual container ID

# echo "Step 5: Stop a container"
# curl -s -X POST "$BASE_URL/api/containers/CONTAINER_ID/stop" \
#   -H "Authorization: Bearer $TOKEN" | jq .
# echo ""

# echo "Step 6: Start a container"
# curl -s -X POST "$BASE_URL/api/containers/CONTAINER_ID/start" \
#   -H "Authorization: Bearer $TOKEN" | jq .
# echo ""

# echo "Step 7: Restart a container"
# curl -s -X POST "$BASE_URL/api/containers/CONTAINER_ID/restart" \
#   -H "Authorization: Bearer $TOKEN" | jq .
# echo ""

echo "=== Demo Completed ==="
