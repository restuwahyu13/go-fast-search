GO = @go
NPM = @npm
NODEMON = @nodemon
DOCKER = @docker
COMPOSE = @docker-compose

#################################
# Application Territory
#################################
.PHONY: install
install:
	${GO} get .
	${GO} mod verify
	${NPM} i nodemon@latest -g

.PHONY: dev
dev:
	${NODEMON} -V -e .go,.env -w . -x go run ./cmd/api --count=1 --race -V --signal SIGTERM

.PHONY: build
build:
	${GO} mod tidy
	${GO} mod verify
	${GO} vet --race -v ./cmd/api
	${GO} build --race -v --ldflags "-r -s -w -extldflags" -o main ./cmd/api

.PHONY: test
test:
	${GO} test -v ./domain/services/**

#################################
# Docker Territory
#################################
.PHONY: upb
upb:
	${DOCKER} build -t go-api:latest --compress .

.PHONY: up
up:
	${COMPOSE} up -d --remove-orphans --no-deps --build

.PHONY: down
down:
	${COMPOSE} down