#!/bin/bash
set -e

# Load .env
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo ".env file not found"
  exit 1
fi

# Validate required variables
: "${DATABASE_HOST:?Missing DATABASE_HOST}"
: "${DATABASE_USERNAME:?Missing DATABASE_USERNAME}"
: "${DATABASE_PASSWORD:?Missing DATABASE_PASSWORD}"
: "${DATABASE_NAME:?Missing DATABASE_NAME}"
: "${DATABASE_PORT:?Missing DATABASE_PORT}"

# Build DATABASE_URL (Postgres URI)
DATABASE_URL="postgres://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable&TimeZone=Asia/Jakarta"

export DATABASE_URL

# Run migrate
migrate -database "$DATABASE_URL" -path ./migrations "$@"