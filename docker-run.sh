#!/bin/bash

set -e

docker build -f processing/Dockerfile -t kashboard-processing:latest .
docker build -f grafana/Dockerfile -t kashboard-grafana:latest .
docker-compose up -d postgres
docker-compose up -d kaspad

sleep 10s

docker-compose up -d processing

sleep 10s

docker-compose up -d grafana

docker-compose logs -f
