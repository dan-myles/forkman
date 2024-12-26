#!/bin/bash

# Enable error handling
set -e

echo "Starting deployment process..."

# Step 2: Clean up old Docker images named 'forkman'
echo "Cleaning up old Docker images..."
docker images -q forkman | xargs -r docker rmi -f

# Step 3: Build a new Docker image named 'forkman'
echo "Building new Docker image..."
docker build -t forkman .

# Step 4: Start the service with Docker Compose
echo "Stopping any existing containers..."
docker compose down


echo "----------------------------------"
echo "SUCCESFULLY BUILT DOCKER IMAGE"
echo "Please use docker-compose to start the service"
echo "----------------------------------"
