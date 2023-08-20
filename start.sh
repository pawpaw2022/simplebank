#!/bin/sh

set -e 

echo "Run db migrations"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Run app"
exec "$@" 

