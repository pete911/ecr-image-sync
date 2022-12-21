NAME := ecr-image-sync
IMAGE := pete911/${NAME}
VERSION ?= dev

test:
	go test ./...
.PHONY:test

build: test
	go build -ldflags "-X main.Version=${VERSION}"
.PHONY:build
