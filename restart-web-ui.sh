#!/bin/bash
set -e

echo "Stopping any running Next.js process..."
pkill -f "next dev" || true

echo "Installing dependencies and rebuilding CSS..."
./install-web-deps.sh

echo "Starting Next.js development server..."
cd web-ui
npm run dev 