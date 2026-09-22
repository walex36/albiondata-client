#!/usr/bin/env bash

set -eo pipefail

TARGET_ARCH="${ARCH:-amd64}"
rm -f albiondata-client
rm -f albiondata-client.gz
rm -f "update-darwin-$TARGET_ARCH.gz"
rm -f "albiondata-client-$TARGET_ARCH-mac.zip"

# Build frontend if dist does not already exist
if [ ! -f "frontend/dist/index.html" ]; then
    (cd frontend && if [ ! -d "node_modules" ]; then npm ci; fi && npm run build)
fi

# Native macOS build. By default builds amd64 (for GitHub CI releases).
# Can specify ARCH=arm64 (e.g. ARCH=arm64 ./scripts/build-darwin.sh) for Apple Silicon.
export CGO_ENABLED=1
TARGET_ARCH="${ARCH:-amd64}"
export GOARCH="$TARGET_ARCH"

if [ "$TARGET_ARCH" = "arm64" ]; then
    export CC="clang -arch arm64"
    export CGO_LDFLAGS="-arch arm64"
else
    export CC="clang -arch x86_64"
    export CGO_LDFLAGS="-arch x86_64"
fi

VERSION="${GITHUB_REF_NAME:-dev}"
go build -ldflags "-s -w -X main.version=$VERSION" -o albiondata-client albiondata-client.go

gzip -k9 albiondata-client
mv albiondata-client.gz "update-darwin-$TARGET_ARCH.gz"

# Zipped folder with a run.command file that runs the client under sudo
TEMP="albiondata-client"
ZIPNAME="albiondata-client-$TARGET_ARCH-mac.zip"
rm -rfv ./scripts/$TEMP
rm -rfv ./$ZIPNAME
mkdir -v ./scripts/$TEMP
cp -v albiondata-client ./scripts/$TEMP/albiondata-client-executable
cp -v ./scripts/run.command ./scripts/$TEMP/run.command
chmod -v 777 ./scripts/$TEMP/*
(cd scripts && zip -v ../$ZIPNAME -r ./"$TEMP")
