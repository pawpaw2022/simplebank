#!/bin/sh

set -e 

echo "Run db migrations"
source /app/app.env
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Run app"
exec "$@" 
