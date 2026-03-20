# BacktesterApp

A web application that allows users to backtest MetaTrader 4 Expert Advisors (.ex4 files) against custom historical data (.csv files). The system is designed to be packaged into a single, fully self-contained Docker image.

## Architecture

This project uses a multi-stage Docker build resulting in a single-container architecture:
- **Frontend**: Svelte 5 with Tailwind CSS (built in Stage 1).
- **Backend**: Golang serving as the API and execution coordinator (built in Stage 2).
- **Execution Environment**: Ubuntu base image containing `wine`, `xvfb`, and the MT4 portable directory to run the backtest headlessly.

The Go backend serves the Svelte UI, accepts `.ex4` and `.csv` files via API upload, builds `config.ini`, executes the MT4 terminal, and returns a JSON-parsed representation of the `.htm` or XML report.

## Directory Structure

- `/frontend` - Svelte 5 user interface.
- `/backend` - Golang server, API handlers, and report parser.
- `/data/experts` - Temporary storage for uploaded `.ex4` files.
- `/data/history` - Temporary storage for uploaded `.csv` history files.

## Getting Started

### 1. Local Environment Bootstrap
Ensure you have the required dependencies (Docker, Go, Node.js/npm, Git) by running the interactive bootstrap script:
- **Linux / macOS**: `./setup.sh`
- **Windows**: `.\setup.ps1`

### 2. Coming Soon
Instructions to build and run the multi-stage Dockerfile will be added here once implemented.
