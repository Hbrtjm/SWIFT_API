#!/bin/bash

# Check if Docker is installed
docker --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
    printf "Docker is not installed. Please install Docker to run this script.\n"
    printf "You can install it by running:\n"
    printf ' $\tsudo apt-get update\n $\tsudo apt-get install ca-certificates curl\n $\tsudo install -m 0755 -d /etc/apt/keyrings\n $\tsudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc\n $\tsudo chmod a+r /etc/apt/keyrings/docker.asc\n $\techo \"deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \$(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \ sudo tee /etc/apt/sources.list.d/docker.list > /dev/null\n $\tsudo apt-get update'
    printf "or by following the instructions at https://docs.docker.com/engine/install/ubuntu/\n"
    exit 1
fi

# Check if Docker Compose is installed
docker compose ps > /dev/null 2>&1 # Any command just to check if docker compose is installed
if [ $? -ne 0 ]; then
    printf "Docker Compose is not installed. Please install Docker Compose to run this script.\n"
    printf "You can install it by running:\n"
    printf " $\tsudo apt-get update\n $\tsudo apt-get install docker-compose-plugin"
    printf "or by following the instructions at https://docs.docker.com/compose/install/linux/\n"
    exit 2
fi

# INFORMATION - If the application is being run as a part of the github actions job, the cleanup will not fail, just continue to deployment

# Clean up any existing containers from previous runs
echo "Cleaning up existing containers..."
docker compose down -v --remove-orphans 2>/dev/null || true

# If it was built as a part of the test, remove any test containers taking up the port 27017
containers=$(docker ps -a -q --filter "name=mongo" 2>/dev/null) || true
if [ -n "$containers" ]; then
  echo "Stopping and removing existing MongoDB containers..."
  docker stop $containers || true
  docker rm $containers || true
fi

# Start the application
echo "Starting application..."
cd ../..
pwd
docker compose up -d --build

echo "Deployment complete. Containers are running:"
docker ps
