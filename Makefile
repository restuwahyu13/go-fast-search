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

.PHONY: adev
adev:
	${NODEMON} -V -e .go,.env -w . -x go run ./cmd/api --count=1 --race -V --signal SIGTERM

.PHONY: abuild
abuild:
	${GO} mod tidy
	${GO} mod verify
	${GO} vet --race -v ./cmd/api
	${GO} build --race -v --ldflags "-r -s -w -extldflags" -o api ./cmd/api

.PHONY: wdev
wdev:
	${NODEMON} -V -e .go,.env -w . -x go run ./cmd/worker --count=1 --race -V --signal SIGTERM

.PHONY: wbuild
wbuild:
	${GO} mod tidy
	${GO} mod verify
	${GO} vet --race -v ./cmd/api
	${GO} build --race -v --ldflags "-r -s -w -extldflags" -o worker ./cmd/worker


.PHONY: test
test:
	${GO} test -v ./domain/services/**

#################################
# Docker Territory
#################################
.PHONY: upb
upb:
	${DOCKER} build -t go-fast-search:latest --compress .

.PHONY: up
up:
	${COMPOSE} up -d --remove-orphans --no-deps

.PHONY: down
down:
	${COMPOSE} down