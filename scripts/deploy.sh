#!/bin/bash

if wevtutil > /dev/null 2>&1; then
    pritnf "You are running on Windows.\n"
    printf "Please use run.ps1 or run.bat to run the application.\n"
    exit 1
fi

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
    exit 1
fi

docker compose up -d --build