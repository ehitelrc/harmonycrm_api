#!/bin/bash

 

echo "Starting Docker build..."

docker buildx build --platform linux/amd64 \
  -t ehitelrc/harmony_service:latest \
  --push \
  .
