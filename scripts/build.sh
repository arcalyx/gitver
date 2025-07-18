#!/bin/bash
set -e

# Parameters
GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}
VERSION=${VERSION:-"1.0.0"}

echo "Building for $GOOS/$GOARCH with version $VERSION"

# Set output filename
OUTPUT="gitver"
if [ "$GOOS" = "windows" ]; then
  OUTPUT="gitver.exe"
fi

# Create dist directory first
mkdir -p dist

# Build the binary
go build -v -ldflags "-X github.com/arcalyx/gitver/cmd.Version=$VERSION" -o "dist/${OUTPUT}"

# Package the binary
cd dist
if [ "$GOOS" = "windows" ]; then
  zip -r "gitver_${VERSION}_${GOOS}_${GOARCH}.zip" "${OUTPUT}"
  echo "Created Windows zip archive"
else
  tar -czvf "gitver_${VERSION}_${GOOS}_${GOARCH}.tar.gz" "${OUTPUT}"
  echo "Created tar.gz archive"
fi
cd ..

# List the contents of the dist directory
echo "Contents of dist directory:"
ls -la dist/
