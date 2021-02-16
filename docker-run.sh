#!/bin/bash

set -e

docker build -f processing/Dockerfile -t kashboard-processing:latest .
docker-compose up -d postgres
docker-compose up -d kaspad

sleep 10s

docker-compose up -d processing

docker-compose logs -f
