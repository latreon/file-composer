#!/bin/bash
set -e

echo "Cleaning up project..."

# Remove empty directories
echo "Removing empty directories and development artifacts..."
rm -rf web-ui/.next

# Clean temporary files
echo "Cleaning up temporary files..."
rm -rf uploads/*
rm -rf compressed/*

# Keep the directories but remove contents
touch uploads/.keep
touch compressed/.keep

echo "Cleanup complete!" 