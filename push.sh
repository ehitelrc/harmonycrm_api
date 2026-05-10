#!/bin/bash

# Validate that a commit message was provided
if [ -z "$1" ]; then
    echo "❌ Error: Faltó proporcionar un mensaje de commit."
    echo "Uso: ./push.sh \"Tu mensaje de commit aquí\""
    exit 1
fi

COMMIT_MESSAGE="$1"

# Check if there are any changes to commit
if [[ -z $(git status -s) ]]; then
    echo "⚠️ No hay cambios pendientes para hacer commit."
    exit 0
fi

# Git operations
echo "📦 Agregando cambios..."
git add .

echo "📝 Haciendo commit: '$COMMIT_MESSAGE'"
git commit -m "$COMMIT_MESSAGE"

echo "🚀 Enviando al repositorio remoto..."
git push

if [ $? -eq 0 ]; then
    echo "✅ Push exitoso en GitHub."
else
    echo "❌ Error al hacer push."
fi
