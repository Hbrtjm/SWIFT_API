@echo off
SETLOCAL EnableDelayedExpansion

REM Check if Docker is installed
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker is not installed. Please install Docker to run this script.
    echo You can install Docker Desktop for Windows by visiting:
    echo     https://docs.docker.com/desktop/install/windows-install/
    echo After installation, ensure Docker Desktop is running.
    exit /b 1
)

REM Check if Docker Compose is installed
docker compose version >nul 2>&1
if %errorlevel% neq 0 (
    echo Docker Compose is not installed. Please install Docker Compose to run this script.
    echo Docker Compose is included with Docker Desktop for Windows.
    echo Make sure you're using a recent version of Docker Desktop that includes Compose V2.
    echo More information at:
    echo     https://docs.docker.com/compose/
    exit /b 1
)

REM Clean up any existing containers
echo Cleaning up existing containers...
docker compose down -v --remove-orphans >nul 2>&1

REM Remove any existing MongoDB containers that might be using port 27017
for /f "tokens=*" %%i in ('docker ps -a -q --filter "name=mongo"') do (
    echo Stopping and removing existing MongoDB container: %%i
    docker stop %%i >nul 2>&1
    docker rm %%i >nul 2>&1
)

REM Start the application
echo Starting application...
docker compose up -d --build

echo Deployment complete. Containers are running:
docker ps

ENDLOCAL
