@echo off
REM Script to create a new release

REM Check if version is provided
if "%~1"=="" (
    echo Error: No version tag provided
    echo Usage: release.bat v1.0.2
    exit /b 1
)

set VERSION=%~1

REM Validate version format (basic check)
echo %VERSION% | findstr /r "^v[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*$" > nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Version tag must be in format v1.0.2
    exit /b 1
)

echo Creating release %VERSION%...

REM Create a tag
git tag %VERSION%

REM Push the tag to trigger the GitHub Actions workflow
git push origin %VERSION%

echo Release %VERSION% created and pushed to GitHub.
echo GitHub Actions will automatically build and release the binaries.
echo Check the status at: https://github.com/georgeglarson/venice-ai-for-brave-leo/actions