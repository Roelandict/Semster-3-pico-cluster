#!/bin/bash

# Build script for multi-platform compilation
# Supports: linux/amd64 (Windows/Linux), linux/arm64 (Raspberry Pi)

set -e

BINARY_NAME="sensor-verwerker"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")

echo "Building sensor-verwerker v${VERSION}"

# Default target
TARGET_ARCH="${1:-amd64}"

case $TARGET_ARCH in
    amd64)
        echo "Building for AMD64 (Windows/Linux x86_64)..."
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BINARY_NAME}-amd64 -ldflags="-X main.Version=${VERSION}" .
        echo "✓ Binary: ${BINARY_NAME}-amd64"
        ;;
    arm64)
        echo "Building for ARM64 (Raspberry Pi)..."
        GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ${BINARY_NAME}-arm64 -ldflags="-X main.Version=${VERSION}" .
        echo "✓ Binary: ${BINARY_NAME}-arm64"
        ;;
    all)
        echo "Building for all platforms..."
        GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BINARY_NAME}-amd64 -ldflags="-X main.Version=${VERSION}" .
        echo "✓ Binary: ${BINARY_NAME}-amd64"
        GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ${BINARY_NAME}-arm64 -ldflags="-X main.Version=${VERSION}" .
        echo "✓ Binary: ${BINARY_NAME}-arm64"
        ;;
    docker)
        echo "Building Docker image for multi-platform (amd64 + arm64)..."
        docker buildx build --platform linux/amd64,linux/arm64 -t sensor-verwerker:${VERSION} -t sensor-verwerker:latest --push .
        echo "✓ Docker image pushed: sensor-verwerker:${VERSION}"
        ;;
    *)
        echo "Usage: ./build.sh [amd64|arm64|all|docker]"
        echo "  amd64  - Build for AMD64 (default)"
        echo "  arm64  - Build for ARM64 (Raspberry Pi)"
        echo "  all    - Build for all platforms"
        echo "  docker - Build and push Docker image for both platforms"
        exit 1
        ;;
esac

echo "Build complete!"
