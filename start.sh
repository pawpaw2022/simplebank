#!/bin/sh

set -e 

echo "Run db migrations"
ls -la /app
source /app/app.env

echo "DB_SOURCE: $DB_SOURCE"
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Run app"
exec "$@" 
