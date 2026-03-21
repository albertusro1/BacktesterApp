#!/bin/bash
set -e

echo "=== Backtester Web Application ==="

# Initialize the Wine prefix headlessly to prevent setup prompts during API calls
echo "Initializing Wine prefix..."
export WINEPREFIX=/root/.wine
xvfb-run -a wineboot --init

echo "Starting Golang Backend..."
cd /app

# Ensure we have the port set correctly
export PORT=3000

# Execute server
exec ./backendapp
