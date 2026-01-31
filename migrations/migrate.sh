#!/bin/bash

export DATABASE_URL="postgres://postgres:root@localhost:5432/order_service?sslmode=disable"
migrate -database "$DATABASE_URL" -path ./migrations "$@"