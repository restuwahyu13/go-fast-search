#!/bin/bash

cd ./apps/fe || exit 1
docker build -t web:latest --compress .
cd ..

wait

cd ./be || exit 1
docker build -f external/deployments/docker/Dockerfile -t api:latest --compress .

sleep 3s

echo "All docker build process completed"
cd ..