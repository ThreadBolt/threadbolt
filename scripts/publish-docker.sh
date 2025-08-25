#!/usr/bin/env bash
set -euo pipefail

# Check for version argument
VERSION=$1
if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker and try again."
    exit 1
fi

# Check if dist/threadbolt-linux-amd64 exists
if [ ! -f "dist/threadbolt-linux-amd64" ]; then
    echo "❌ Binary dist/threadbolt-linux-amd64 not found. Run ./scripts/prepare-release.sh $VERSION first."
    exit 1
fi

# Check if Dockerfile exists
if [ ! -f "Dockerfile" ]; then
    echo "❌ Dockerfile not found. Please create a Dockerfile in the repository root."
    exit 1
fi

# Build Docker image
echo "🐳 Building Docker image for ThreadBolt CLI v$VERSION..."
DOCKER_USERNAME="haiderzamanyzi" # Replace with your Docker Hub username
docker build -t $DOCKER_USERNAME/threadbolt:$VERSION .

# Log in to Docker Hub
echo "🔐 Logging in to Docker Hub..."
docker login || { echo "❌ Docker login failed. Please check your credentials."; exit 1; }

# Push Docker image
echo "🚀 Pushing Docker image to Docker Hub..."
docker push $DOCKER_USERNAME/threadbolt:$VERSION

# Tag and push latest
echo "🏷️ Tagging and pushing latest image..."
docker tag $DOCKER_USERNAME/threadbolt:$VERSION $DOCKER_USERNAME/threadbolt:latest
docker push $DOCKER_USERNAME/threadbolt:latest

echo "✅ Docker image published as $DOCKER_USERNAME/threadbolt:$VERSION and $DOCKER_USERNAME/threadbolt:latest"