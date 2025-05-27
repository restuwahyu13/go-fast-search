NPM = @npm
DOCKER = @docker
COMPOSE = @docker-compose
DOCKERFILE_FE = $(realpath ./apps/fe)
DOCKERFILE_BE = $(realpath ./apps/be/external/deployments/docker)

#################################
# Application Territory
#################################
.PHONY: dev
dev:
	${NPM} run dev

.PHONY: install
install:
	${NPM} run install

.PHONY: worker
worker:
	${NPM} run worker

.PHONY: scheduler
scheduler:
	${NPM} run scheduler

.PHONY: build
build:
	./app-build.sh

#################################
# Docker Territory
#################################
.PHONY: upb
upb:
	./docker-build.sh

.PHONY: up
up:
	${COMPOSE} up -d --remove-orphans --no-deps

.PHONY: down
down:
	${COMPOSE} down