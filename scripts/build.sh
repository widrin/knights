#!/bin/bash

set -e

echo "Building Knights Game Server..."

# Build server
go build -o bin/server cmd/server/main.go

echo "Build complete! Binary: bin/server"
