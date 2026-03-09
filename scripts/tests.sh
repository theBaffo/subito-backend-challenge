#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="subito-backend-challenge-test"

docker build -t "$IMAGE_NAME" --target builder .
docker run --rm "$IMAGE_NAME" npm test
