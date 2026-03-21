# Stage 1: Build the Svelte Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci || npm install
COPY frontend .
RUN npm run build

# Stage 2: Build the Golang Backend
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend .
RUN GOOS=linux GOARCH=amd64 go build -o backendapp .

# Stage 3: Ubuntu Base Image for Execution Environment
FROM ubuntu:22.04
ENV DEBIAN_FRONTEND=noninteractive

# Add architecture i386 for Wine and install xvfb
RUN dpkg --add-architecture i386 && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    wine32 wine xvfb ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Create necessary directories
RUN mkdir -p /data/MT4 /data/experts /data/history

# Copy Frontend build
COPY --from=frontend-builder /app/dist /app/frontend/dist

# Copy Backend binary
COPY --from=backend-builder /app/backendapp /app/backendapp

# Copy entrypoint
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 3000

ENTRYPOINT ["/app/entrypoint.sh"]
