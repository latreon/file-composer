#!/bin/bash

# Exit script if any command fails
set -e

# Build the API server
echo "Building API server..."
go build -o build/api-server ./cmd/api

# Install frontend dependencies (if needed)
if [ ! -d "web-ui/node_modules" ]; then
  echo "Installing frontend dependencies..."
  cd web-ui
  npm install
  cd ..
fi

# Run the API server in the background
echo "Starting API server..."
./build/api-server &
API_PID=$!

# Wait a moment for the API server to start
sleep 2

# Run the Next.js frontend
echo "Starting Next.js frontend..."
cd web-ui
npm run dev &
NEXT_PID=$!

# Function to handle script termination
cleanup() {
  echo "Shutting down services..."
  kill $NEXT_PID
  kill $API_PID
  exit 0
}

# Register the cleanup function for script termination
trap cleanup SIGINT SIGTERM

echo ""
echo "Services running:"
echo "API server: http://localhost:8080"
echo "Frontend: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop all services"

# Wait for user to press Ctrl+C
wait 