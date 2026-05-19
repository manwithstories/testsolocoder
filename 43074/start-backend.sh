#!/bin/bash

echo "Starting Book Library Backend..."
cd "$(dirname "$0")/backend"
go run cmd/server/main.go
