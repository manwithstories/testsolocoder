#!/bin/bash

# ============================================================
# Recruitment Platform - Setup Script
# ============================================================

set -e

echo "============================================="
echo "  Recruitment Platform Setup"
echo "============================================="

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

# Check for Node.js
if ! command -v node &> /dev/null; then
    echo "Error: Node.js is not installed. Please install Node.js first."
    exit 1
fi

# Check for PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "Warning: PostgreSQL client not found. Make sure PostgreSQL is installed."
fi

echo ""
echo "Setting up Backend..."
cd backend

# Initialize Go modules
go mod tidy

# Create uploads directory
mkdir -p uploads/resumes

echo "Backend setup complete!"
echo ""

echo "Setting up Frontend..."
cd ../frontend

# Install dependencies
npm install

echo "Frontend setup complete!"
echo ""

echo "============================================="
echo "  Setup Complete!"
echo "============================================="
echo ""
echo "To start the application:"
echo ""
echo "1. Start PostgreSQL and create database:"
echo "   createdb recruitment_db"
echo ""
echo "2. Start Backend (in /backend directory):"
echo "   go run main.go"
echo ""
echo "3. Start Frontend (in /frontend directory):"
echo "   npm run dev"
echo ""
echo "4. Access the application at http://localhost:5173"
echo ""
