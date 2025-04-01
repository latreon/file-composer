#!/bin/bash
set -e

echo "Installing web UI dependencies..."
cd web-ui

# Install all dependencies from package.json
npm install

# Ensure Tailwind CSS and related packages are installed
npm install --save tailwindcss@latest postcss@latest autoprefixer@latest

# Initialize Tailwind if not already done (won't overwrite existing config)
if [ ! -f "tailwind.config.js" ]; then
  npx tailwindcss init -p
fi

echo "Dependencies installed successfully!"
echo "You can now run the application with ./run.sh" 