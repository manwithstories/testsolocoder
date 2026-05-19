#!/bin/bash

set -e

echo "=== Freelancer Management System ==="
echo ""

cd "$(dirname "$0")/.."

echo "1. Starting backend server..."
cd backend
go mod download
go run cmd/server/main.go &
BACKEND_PID=$!
echo "   Backend running on http://localhost:8080 (PID: $BACKEND_PID)"

sleep 2

echo ""
echo "2. Starting frontend dev server..."
cd ../frontend
if [ ! -d "node_modules" ]; then
    echo "   Installing frontend dependencies..."
    npm install
fi
npm run dev &
FRONTEND_PID=$!
echo "   Frontend running on http://localhost:3000 (PID: $FRONTEND_PID)"

echo ""
echo "=== Servers started ==="
echo "Backend: http://localhost:8080"
echo "Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all servers"

trap "echo 'Stopping servers...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" SIGINT SIGTERM

wait
