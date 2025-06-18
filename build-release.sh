#!/bin/bash
# Build script for creating release binaries for multiple platforms

set -e

VERSION="1.0.0"
BINARY_NAME="orgmarks"

echo "Building orgmarks version $VERSION for multiple platforms..."

# Create release directory
mkdir -p release

# Linux AMD64
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o "release/${BINARY_NAME}-${VERSION}-linux-amd64" -ldflags="-s -w"

# Linux ARM64
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -o "release/${BINARY_NAME}-${VERSION}-linux-arm64" -ldflags="-s -w"

# macOS AMD64 (Intel)
echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -o "release/${BINARY_NAME}-${VERSION}-darwin-amd64" -ldflags="-s -w"

# macOS ARM64 (Apple Silicon)
echo "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -o "release/${BINARY_NAME}-${VERSION}-darwin-arm64" -ldflags="-s -w"

# Windows AMD64
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -o "release/${BINARY_NAME}-${VERSION}-windows-amd64.exe" -ldflags="-s -w"

echo ""
echo "Build complete! Binaries are in the release/ directory:"
ls -lh release/

echo ""
echo "To create archives:"
echo "  cd release"
echo "  tar czf ${BINARY_NAME}-${VERSION}-linux-amd64.tar.gz ${BINARY_NAME}-${VERSION}-linux-amd64"
echo "  tar czf ${BINARY_NAME}-${VERSION}-linux-arm64.tar.gz ${BINARY_NAME}-${VERSION}-linux-arm64"
echo "  tar czf ${BINARY_NAME}-${VERSION}-darwin-amd64.tar.gz ${BINARY_NAME}-${VERSION}-darwin-amd64"
echo "  tar czf ${BINARY_NAME}-${VERSION}-darwin-arm64.tar.gz ${BINARY_NAME}-${VERSION}-darwin-arm64"
echo "  zip ${BINARY_NAME}-${VERSION}-windows-amd64.zip ${BINARY_NAME}-${VERSION}-windows-amd64.exe"
