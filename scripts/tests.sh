#!/bin/sh
set -e

echo "Building test image..."
docker build --target builder -t subito-backend-challenge-test .

echo "Running test suite..."
docker run --rm subito-backend-challenge-test \
  go test ./... -v -count=1
