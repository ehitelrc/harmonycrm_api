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

echo "============================================="
echo "🛠️ Realizando pruebas de seguridad..."
echo "============================================="

echo "1. Verificando compilación (go build)..."
if ! go build -o /tmp/harmony_test_bin .; then
    echo "❌ ERROR: La aplicación no compila. Abortando despliegue."
    exit 1
fi
echo "✅ Compilación exitosa."

echo "2. Prueba de arranque en frío (Dry Run)..."
# Ejecutamos el binario en segundo plano para ver si hay un panic de inicialización (como el de Gin)
/tmp/harmony_test_bin &
APP_PID=$!

# Esperamos 3 segundos para darle tiempo a que inicialice rutas y base de datos
sleep 3

# Comprobamos si el proceso sigue vivo
if ! kill -0 $APP_PID 2>/dev/null; then
    echo "❌ ERROR CATASTRÓFICO: La aplicación hizo 'panic' al arrancar (probablemente un conflicto de rutas o DB)."
    echo "Abortando despliegue para proteger producción."
    rm /tmp/harmony_test_bin
    exit 1
fi

# Si sigue vivo, lo matamos de forma segura porque pasó la prueba
kill $APP_PID
rm /tmp/harmony_test_bin
echo "✅ Prueba de arranque exitosa (Sin panics iniciales)."
echo "============================================="

echo "Starting Docker build..."

docker buildx build --platform linux/amd64 \
  -t ehitelrc/harmony_service:latest \
  --push \
  .
