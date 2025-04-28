# Check if running on Windows
If ($OsType -eq "Linux")
{
    Write-Host "This script is intended to be run on Windows. Please run it on a Windows machine."
    Exit 1
}

# Check if Docker is installed
If (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Docker is not installed. Please install Docker to run this script: https://docs.docker.com/desktop/windows/install/"
    Exit 1
}

# Check if Docker Compose is installed
Try {
    docker compose ps | Out-Null
} Catch {
    Write-Host "Docker Compose is not installed. Please install Docker Compose to run this script: https://docs.docker.com/desktop/setup/install/windows-install/"
    Exit 1
}

# Create Docker network (replace with actual driver and subnet if needed)
docker network create --gateway 177.199.0.1 --driver bridge swift_api --subnet 177.199.0.0/24 | Out-Null

# Build and start the Docker containers
docker compose up -d --build