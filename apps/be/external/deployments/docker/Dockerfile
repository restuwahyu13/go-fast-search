# ======================
#  GO STAGE
# ======================
FROM golang:latest AS builder

LABEL name="Restu Wahyu Saputra"
LABEL email="restuwahyu13@gmail.com"
LABEL linkedin="https://www.linkedin.com/in/restuwahyu13"
LABEL gitHub="https://github.com/restuwahyu13"

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download; \
    go mod verify

COPY ./ .
RUN go vet ./cmd/*; \
    go build -v -ldflags="-s -w" -o api ./cmd/api; \
    go build -v -ldflags="-s -w" -o worker ./cmd/worker; \
    go build -v -ldflags="-s -w" -o scheduler ./cmd/scheduler

# ======================
#  ALPINE STAGE
# ======================
FROM alpine:latest
WORKDIR /usr/src/app

COPY --from=builder /app/api /app/worker /app/scheduler ./

RUN apk update; \
    apk -u list; \
    apk upgrade; \
    apk add --no-cache tzdata ca-certificates