#!/bin/bash
# Build script for Leo Venice.AI Configuration Tool (Windows only)

# Create a build directory
mkdir -p build

echo "Building Leo Venice.AI Configuration Tool for Windows..."

# Build for Windows (64-bit)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o build/leo_venice_config.exe leo_venice_config.go

echo "Build process completed!"
echo "Executable is available in the 'build' directory: build/leo_venice_config.exe"