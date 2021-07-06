#!/bin/bash

set -e

echo "Waiting for postgres..."

while ! nc -z db 5432; do
  sleep 0.1
done

echo "PostgreSQL started"

migrate -source file://db/migrations -database postgres://postgres:password@db:5432/users?sslmode=disable up

