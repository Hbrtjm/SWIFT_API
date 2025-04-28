@echo off

REM Check if running on Windows
ver | findstr /i "Windows" >nul
if %errorlevel% neq 0 (
    echo This script is intended to be run on Windows. Please run it on a Windows machine.
    exit /b 1
)

REM Check if Docker is installed
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not installed. Please install Docker to run this script: https://docs.docker.com/desktop/windows/install/
    exit /b 1
)

REM Check if Docker Compose is installed
docker compose version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker Compose is not installed. Please install Docker Compose to run this script: https://docs.docker.com/desktop/setup/install/windows-install/
    exit /b 1
)

REM Create Docker network (replace with actual driver and subnet if needed)
docker network create --gateway 177.199.0.1 --driver bridge swift_api --subnet 177.199.0.0/24 >nul 2>&1

REM Build and start the Docker containers
docker compose up -d --build