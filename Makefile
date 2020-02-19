SHELL := /bin/bash
GIT_COMMIT := $(shell git rev-list -1 HEAD)
PROJECT_NAME := go-graceful

PHONY: run run-git-info run-prod build deps
deps:
	go mod tidy && go mod vendor

run:
	APP_GIT_BUILD_VERSION=$(GIT_COMMIT) go run -mod=vendor cmd/std/main.go

run-git-info:
	APP_GIT_BUILD_VERSION=$(GIT_COMMIT) mvn go run -mod=vendor cmd/std/main.go

build:
	go build -mod=vendor -o out/goserver-$(GIT_COMMMIT) cmd/std/main.go

run-prod: build
	APP_GIT_BUILD_VERSION=$(GIT_COMMIT) ./out/goserver-$(GIT_COMMIT) -port=8082


