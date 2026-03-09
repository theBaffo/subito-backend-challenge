#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="subito-backend-challenge"

docker build -t "$IMAGE_NAME" .
docker run --rm -p 3000:3000 "$IMAGE_NAME"
