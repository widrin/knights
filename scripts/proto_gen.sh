#!/bin/bash

set -e

echo "Generating protobuf code..."

# Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative \
    pkg/proto/*.proto

echo "Protobuf generation complete!"
