#!/usr/bin/env bash

echo "Building for Linux..."
GOOS=linux go build -o ./build/sample-linux
echo "Pushing..."
cf push -f manifest.yml