# Check if Docker is installed
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Docker is not installed. Please install Docker to run this script."
    Write-Host "You can install Docker Desktop for Windows by visiting:"
    Write-Host "    https://docs.docker.com/desktop/install/windows-install/"
    Write-Host "After installation, ensure Docker Desktop is running."
    exit 1
}

# Check if Docker Compose is installed
try {
    docker compose version | Out-Null
} catch {
    Write-Host "Docker Compose is not installed. Please install Docker Compose to run this script."
    Write-Host "Docker Compose is included with Docker Desktop for Windows."
    Write-Host "Make sure you're using a recent version of Docker Desktop that includes Compose V2."
    Write-Host "You can verify this by running: docker compose version"
    Write-Host "More info at:"
    Write-Host "    https://docs.docker.com/compose/"
    exit 1
}

# Clean up existing containers
Write-Host "Cleaning up existing containers..."
try {
    docker compose down -v --remove-orphans | Out-Null
} catch {
    # Ignore any errors
}

# If it was built as a part of the test, remove any test containers taking up the port 27017
$containers = docker ps -a -q --filter "name=mongo"
if ($containers) {
    Write-Host "Stopping and removing existing MongoDB containers..."
    foreach ($container in $containers) {
        try {
            docker stop $container | Out-Null
            docker rm $container | Out-Null
        } catch {
            # Ignore errors
        }
    }
}

# Start the application
Write-Host "Starting application..."
docker compose up -d --build

Write-Host "Deployment complete. Containers are running:"
docker ps
