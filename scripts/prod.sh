#!/bin/bash

# Enable error handling
set -e

echo "Starting deployment process..."

# Step 1: Pull the latest changes from Git
echo "Pulling latest changes from Git..."
git pull origin main || { echo "Failed to pull from Git. Exiting..."; exit 1; }

# Step 2: Clean up old Docker images named 'forkman'
echo "Cleaning up old Docker images..."
docker images -q forkman | xargs -r docker rmi -f

# Step 3: Build a new Docker image named 'forkman'
echo "Building new Docker image..."
docker build -t forkman .

# Step 4: Start the service with Docker Compose
echo "Starting service with Docker Compose..."
docker-compose down
docker-compose up -d

# Step 5: Verify the deployment
echo "Verifying the running containers..."
docker ps | grep forkman

echo "Deployment completed successfully!"
