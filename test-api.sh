#!/bin/bash

# API Testing Script for docker-simple-panel
# This script tests all API endpoints

BASE_URL="http://localhost:8080"
CONTENT_TYPE="Content-Type: application/json"

echo "=== Docker Simple Panel API Tests ==="
echo ""

# Test 1: Health Check
echo "1. Testing Health Endpoint..."
curl -s "$BASE_URL/api/health" | jq .
echo ""

# Test 2: Register User
echo "2. Testing User Registration..."
curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "$CONTENT_TYPE" \
  -d '{
    "username": "testuser",
    "password": "testpass123",
    "nickname": "Test User"
  }' | jq .
echo ""

# Test 3: Login
echo "3. Testing User Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "$CONTENT_TYPE" \
  -d '{
    "username": "testuser",
    "password": "testpass123"
  }')
echo "$LOGIN_RESPONSE" | jq .
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token')
echo ""

# Test 4: List Containers (Authenticated)
echo "4. Testing List Containers (Authenticated)..."
curl -s "$BASE_URL/api/containers" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Test 5: Docker Health Check (Authenticated)
echo "5. Testing Docker Health Check (Authenticated)..."
curl -s "$BASE_URL/api/docker/health" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Test 6: Unauthorized Access
echo "6. Testing Unauthorized Access (Should Fail)..."
curl -s "$BASE_URL/api/containers" -w "\nHTTP Status: %{http_code}\n"
echo ""

# Test 7: Invalid Credentials
echo "7. Testing Invalid Credentials (Should Fail)..."
curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "$CONTENT_TYPE" \
  -d '{
    "username": "testuser",
    "password": "wrongpassword"
  }' | jq .
echo ""

echo "=== All Tests Completed ==="
