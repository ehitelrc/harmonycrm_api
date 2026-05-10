#!/bin/bash

VERSION_FILE="version.txt"

# Crear el archivo si no existe
if [ ! -f "$VERSION_FILE" ]; then
    echo "3.50.0" > "$VERSION_FILE"
fi

# Leer la versión actual
CURRENT_VERSION=$(cat "$VERSION_FILE")

# Separar la versión por puntos
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR="${VERSION_PARTS[0]}"
MINOR="${VERSION_PARTS[1]}"
PATCH="${VERSION_PARTS[2]}"

# Incrementar el parche (x.x.+1)
PATCH=$((PATCH + 1))
NEW_VERSION="$MAJOR.$MINOR.$PATCH"

# Guardar la nueva versión
echo "$NEW_VERSION" > "$VERSION_FILE"
echo "⬆️ Versión incrementada a: $NEW_VERSION"

echo "Starting Docker build..."

docker buildx build --platform linux/amd64 \
  -t ehitelrc/harmony_service:latest \
  --push \
  .
