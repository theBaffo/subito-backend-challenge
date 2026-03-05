#!/bin/sh
set -e

echo "Building Docker image..."
docker build -t subito-backend-challenge .

echo "Starting service on http://localhost:8080"
docker run --rm -p 8080:8080 subito-backend-challenge
