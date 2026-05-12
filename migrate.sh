#!/bin/bash

# Configuration (Use env vars in real production environments)
DB_HOST="20.81.232.132"
DB_USER="harmony"
DB_PASS="harmonyValle2025"
DB_NAME="harmony"
DB_PORT="5432"

DSN="host=$DB_HOST user=$DB_USER password=$DB_PASS dbname=$DB_NAME port=$DB_PORT sslmode=disable"

if [ "$1" == "up" ]; then
    echo "🔼 Ejecutando migraciones pendientes..."
    ~/go/bin/goose -dir migrations postgres "$DSN" up
elif [ "$1" == "status" ]; then
    echo "📊 Estado de las migraciones..."
    ~/go/bin/goose -dir migrations postgres "$DSN" status
elif [ "$1" == "create" ] && [ -n "$2" ]; then
    echo "📄 Creando nueva migración: $2"
    ~/go/bin/goose -dir migrations create "$2" sql
else
    echo "Uso:"
    echo "  ./migrate.sh status         (Muestra el estado de la base de datos)"
    echo "  ./migrate.sh create <name>  (Crea una nueva migración vacía)"
    echo "  ./migrate.sh up             (Aplica todas las migraciones nuevas a producción)"
fi
