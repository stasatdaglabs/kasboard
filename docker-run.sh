#!/bin/bash

set -e

docker build -f processing/Dockerfile -t kashboard-processing:latest .
docker build -f web/Dockerfile -t kashboard-web:latest .
docker-compose up -d postgres
docker-compose up -d kaspad

sleep 10s

docker-compose up -d processing
docker-compose up -d web

docker-compose logs -f
