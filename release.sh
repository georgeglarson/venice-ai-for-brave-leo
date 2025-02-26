#!/bin/bash
# Script to create a new release

# Check if version is provided
if [ $# -eq 0 ]; then
    echo "Error: No version tag provided"
    echo "Usage: ./release.sh v1.0.2"
    exit 1
fi

VERSION=$1

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version tag must be in format v1.0.2"
    exit 1
fi

echo "Creating release $VERSION..."

# Create a tag
git tag $VERSION

# Push the tag to trigger the GitHub Actions workflow
git push origin $VERSION

echo "Release $VERSION created and pushed to GitHub."
echo "GitHub Actions will automatically build and release the binaries."
echo "Check the status at: https://github.com/georgeglarson/venice-ai-for-brave-leo/actions"