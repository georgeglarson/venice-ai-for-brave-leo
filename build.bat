@echo off
REM Build script for Leo Venice.AI Configuration Tool

REM Ensure we have the required dependency
go get github.com/google/uuid

REM Create a build directory
if not exist build mkdir build

echo Building Leo Venice.AI Configuration Tool for multiple platforms...

REM Build for Windows
echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -o build\leo_venice_config.exe leo_venice_config.go
if %ERRORLEVEL% EQU 0 (
    echo ✓ Windows build successful: build\leo_venice_config.exe
) else (
    echo ✗ Windows build failed
)

REM Build for macOS
echo Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -o build\leo_venice_config_mac leo_venice_config.go
if %ERRORLEVEL% EQU 0 (
    echo ✓ macOS build successful: build\leo_venice_config_mac
) else (
    echo ✗ macOS build failed
)

REM Build for Linux
echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -o build\leo_venice_config_linux leo_venice_config.go
if %ERRORLEVEL% EQU 0 (
    echo ✓ Linux build successful: build\leo_venice_config_linux
) else (
    echo ✗ Linux build failed
)

echo Build process completed!
echo Executables are available in the 'build' directory:
echo - Windows: build\leo_venice_config.exe
echo - macOS: build\leo_venice_config_mac
echo - Linux: build\leo_venice_config_linux

REM Reset environment variables
set GOOS=
set GOARCH=