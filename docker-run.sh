#!/bin/bash

set -e

# Build the grafana and processing images
docker build -f processing/Dockerfile -t kashboard-processing:latest .
docker build -f grafana/Dockerfile -t kashboard-grafana:latest .

# Start postgres
docker-compose up -d postgres

# Wait for postgres to finish initializing
sleep 10s

# Start processing
docker-compose up -d processing

# Wait for processing to finish initializing
sleep 10s

# Start grafana
docker-compose up -d grafana

# Print logs for all services
docker-compose logs -f
