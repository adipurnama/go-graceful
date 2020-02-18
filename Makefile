SHELL := /bin/bash
GIT_COMMIT := $(shell git rev-list -1 HEAD)
PROJECT_NAME := go-graceful

PHONY: run run-git-info run-prod build
run:
	APP_GIT_BUILD_VERSION=$(GIT_COMMIT) go run cmd/main.go

run-git-info:
	APP_GIT_BUILD_VERSION=$(GIT_COMMIT) mvn go run cmd/main.go

build:
	go build -o out/goserver cmd/main.go

run-prod: build
	APP_GIT_BUILD_VERSION=$(GIT_COMMIT) ./out/goserver -port=8082


