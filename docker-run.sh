#!/bin/bash

set -e

# Verify that all the required environment variables are set
declare -A REQUIRED_VARIABLES
REQUIRED_VARIABLES["POSTGRES_USER"]="${POSTGRES_USER}"
REQUIRED_VARIABLES["POSTGRES_PASSWORD"]="${POSTGRES_PASSWORD}"
REQUIRED_VARIABLES["POSTGRES_DB"]="${POSTGRES_DB}"
REQUIRED_VARIABLES["KASPAD_ADDRESS"]="${KASPAD_ADDRESS}"
REQUIRED_VARIABLES["GRAFANA_ADMIN_USERNAME"]="${GRAFANA_ADMIN_USERNAME}"
REQUIRED_VARIABLES["GRAFANA_ADMIN_PASSWORD"]="${GRAFANA_ADMIN_PASSWORD}"
REQUIRED_VARIABLES["GRAFANA_PORT"]="${GRAFANA_PORT}"

REQUIRED_VARIABLE_NOT_SET=false
for REQUIRED_VARIABLE_NAME in "${!REQUIRED_VARIABLES[@]}"; do
  if [ -z "${REQUIRED_VARIABLES[$REQUIRED_VARIABLE_NAME]}" ]; then
    echo "${REQUIRED_VARIABLE_NAME} is not set";
    REQUIRED_VARIABLE_NOT_SET=true
    fi
done

if [ true = "${REQUIRED_VARIABLE_NOT_SET}" ]; then
  echo
  echo "The following environment variables are required:"
  for REQUIRED_VARIABLE_NAME in "${!REQUIRED_VARIABLES[@]}"; do
    echo "${REQUIRED_VARIABLE_NAME}"
  done
  exit 1
fi

# Build the grafana and processing images
docker build -f processing/Dockerfile -t kasboard-processing:latest .
docker build -f grafana/Dockerfile -t kasboard-grafana:latest .

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
