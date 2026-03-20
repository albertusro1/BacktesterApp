#!/bin/bash
set -e

echo "=== BacktesterApp Local Environment Setup ==="
echo "Checking for required dependencies..."

deps=("docker" "go" "node" "npm" "git")
missing=()

for dep in "${deps[@]}"; do
    if ! command -v "$dep" &> /dev/null; then
        missing+=("$dep")
    fi
done

if [ ${#missing[@]} -ne 0 ]; then
    echo "Missing dependencies detected: ${missing[*]}"
    read -r -p "Please enter the absolute path for the installation directory (or press Enter for default system paths): " INSTALL_DIR
    
    echo "Starting installation for missing dependencies..."
    
    if [ -z "$INSTALL_DIR" ]; then
        echo "Using default system paths. This may require sudo."
        
        if command -v apt-get &> /dev/null; then
            sudo apt-get update
            for dep in "${missing[@]}"; do
                if [ "$dep" == "docker" ]; then
                    sudo apt-get install -y docker.io
                elif [ "$dep" == "go" ]; then
                    sudo apt-get install -y golang-go
                elif [ "$dep" == "node" ] || [ "$dep" == "npm" ]; then
                    sudo apt-get install -y nodejs npm
                elif [ "$dep" == "git" ]; then
                    sudo apt-get install -y git
                fi
            done
        elif command -v brew &> /dev/null; then
            for dep in "${missing[@]}"; do
                brew install "$dep"
            done
        else
            echo "Unsupported package manager for automatic default installation. Please install manually."
            exit 1
        fi
    else
        echo "Custom installation to $INSTALL_DIR is requested."
        mkdir -p "$INSTALL_DIR"
        echo "Please note: Fully automated custom-path installation for system-level tools like Docker via bash script implies downloading portable binaries."
        echo "For this setup, extracting to $INSTALL_DIR..."
        
        for dep in "${missing[@]}"; do
            echo "Custom install step for $dep to $INSTALL_DIR (Placeholder - Add your specific binary links here)"
        done
    fi
else
    echo "All required dependencies (Docker, Golang, Node.js, npm, Git) are installed!"
fi

echo "Local environment bootstrap complete."
