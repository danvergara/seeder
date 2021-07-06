#!/bin/bash

set -e

echo "Waiting for postgres..."

while ! nc -z db 5432; do
  sleep 0.1
done

echo "PostgreSQL started"

echo "Running the migrations against the database"
make migrate

echo "Seeding the database"
go run db/main.go
