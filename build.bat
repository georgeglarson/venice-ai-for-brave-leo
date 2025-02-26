@echo off
REM Build script for Leo Venice.AI Configuration Tool (Windows only)

REM Ensure we have the required dependency
go get github.com/google/uuid

REM Create a build directory
if not exist build mkdir build

echo Building Leo Venice.AI Configuration Tool for Windows...

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

echo Build process completed!
echo Executable is available in the 'build' directory: build\leo_venice_config.exe

REM Reset environment variables
set GOOS=
set GOARCH=