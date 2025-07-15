#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script configuration
VERSION=${1:-"latest"}
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
VCS_REF=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
IMAGE_NAME="guestbook-api"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Guest Book API Docker Build Script   ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Display build information
echo -e "${YELLOW}Build Information:${NC}"
echo -e "  Version: ${GREEN}${VERSION}${NC}"
echo -e "  Build Date: ${GREEN}${BUILD_DATE}${NC}"
echo -e "  VCS Ref: ${GREEN}${VCS_REF}${NC}"
echo -e "  Image Name: ${GREEN}${IMAGE_NAME}${NC}"
echo ""

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Run tests before building
echo -e "${YELLOW}Running tests...${NC}"
if ! go test ./... -short; then
    echo -e "${RED}Tests failed. Aborting build.${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Tests passed${NC}"
echo ""

# Build the Docker image
echo -e "${YELLOW}Building Docker image...${NC}"
docker build \
    --build-arg VERSION="${VERSION}" \
    --build-arg BUILD_DATE="${BUILD_DATE}" \
    --build-arg VCS_REF="${VCS_REF}" \
    --tag "${IMAGE_NAME}:${VERSION}" \
    --tag "${IMAGE_NAME}:latest" \
    .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Docker image built successfully${NC}"
else
    echo -e "${RED}✗ Docker build failed${NC}"
    exit 1
fi

# Display image information
echo ""
echo -e "${YELLOW}Image Information:${NC}"
docker images "${IMAGE_NAME}:${VERSION}" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

# Optional: Test the built image
echo ""
echo -e "${YELLOW}Testing the built image...${NC}"
CONTAINER_ID=$(docker run -d -p 4261:4260 "${IMAGE_NAME}:${VERSION}")

# Wait for container to start
sleep 5

# Test health endpoint
if curl -f http://localhost:4261/health >/dev/null 2>&1; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed${NC}"
fi

# Cleanup test container
docker stop "${CONTAINER_ID}" >/dev/null 2>&1
docker rm "${CONTAINER_ID}" >/dev/null 2>&1

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Build completed successfully!         ${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo -e "  Run: ${BLUE}docker run -p 4260:4260 ${IMAGE_NAME}:${VERSION}${NC}"
echo -e "  Or:  ${BLUE}docker-compose -f docker-compose.prod.yml up${NC}"
echo ""
