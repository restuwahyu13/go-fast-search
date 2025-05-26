#!/bin/bash

cd ./apps/be || exit 1
rm -r ./api ./worker ./scheduler;

$(go mod verify);
$(go vet --race -v ./cmd/api);
$(go vet --race -v ./cmd/worker);
$(go vet --race -v ./cmd/scheduler);

wait

$(npm run build);
$(go build --race -v -ldflags="-s -w" -o api ./cmd/api);
$(go build --race -v -ldflags="-s -w" -o worker ./cmd/worker);
$(go build --race -v -ldflags="-s -w" -o scheduler ./cmd/scheduler);


sleep 3s

echo "All build app process completed"
cd ..