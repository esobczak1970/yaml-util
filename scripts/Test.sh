#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Running go mod tidy..."
go mod tidy

echo "Running go vet..."
go vet ./... || { echo "go vet failed! Fix issues before running tests."; exit 1; }

echo "Running tests for minify, maxify, and verbose..."
go test ./minify ./maxify ./verbose

echo "✅ All tests passed successfully!"
