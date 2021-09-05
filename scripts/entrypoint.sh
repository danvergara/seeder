#!/bin/bash

set -e

echo "Waiting for postgres..."

while ! nc -z db 5432; do
  sleep 0.1
done

echo "PostgreSQL started"

echo "Running the migrations against the database"
make migrate

echo "Running the seeds against the database"
./seeder -p ./example
