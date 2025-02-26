#!/bin/bash
# Build script for Leo Venice.AI Configuration Tool

# Create a build directory
mkdir -p build

echo "Building Leo Venice.AI Configuration Tool for multiple platforms..."

# Build for Windows (64-bit)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o build/leo_venice_config.exe leo_venice_config.go

# Build for macOS (64-bit)
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o build/leo_venice_config_mac leo_venice_config.go

# Build for Linux (64-bit)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o build/leo_venice_config_linux leo_venice_config.go

echo "Build process completed!"
echo "Executables are available in the 'build' directory:"
echo "- Windows: build/leo_venice_config.exe"
echo "- macOS: build/leo_venice_config_mac"
echo "- Linux: build/leo_venice_config_linux"