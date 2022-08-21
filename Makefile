.DEFAULT_GOAL := dev
GOARCH := amd64 # change to arm64 if on mac m1/m2 or surface
PWD := $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

.PHONY: dev
dev: ## test lint cleanup
dev: go-clean generate vet fmt lint test mod-tidy

.PHONY: build build-linux build-win build-mac
build: ## build server binary
build: build-linux
build-linux: generate mod-tidy
	CGO_ENABLED=1 GOOS=linux GOARCH=${GOARCH} go build -ldflags='-w -s -extldflags "-static"' -o server cmd/server/main.go
build-win: generate mod-tidy
	CGO_ENABLED=1 GOOS=windows GOARCH=${GOARCH} go build -ldflags='-w -s -extldflags "-static"' -o server.exe cmd/server/main.go
build-mac: generate mod-tidy
	CGO_ENABLED=1 GOOS=darwin GOARCH=${GOARCH} go build -ldflags='-w -s -extldflags "-static"' -o server cmd/server/main.go

.PHONY: generate
generate: ## go generate
	docker run -v "${PWD}":/src -w /src vektra/mockery --all --keeptree
	go generate ./...

.PHONY: vet
vet: ## go vet
	go vet ./...

.PHONY: fmt
fmt: ## go fmt
	go fmt ./...

.PHONY: lint
lint: ## golangci-lint
	golangci-lint run

.PHONY: test
test: # test
test:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: mod-tidy
mod-tidy: ## go mod tidy
	go mod tidy

.PHONY: go-clean
go-clean: ## go clean build, test and modules caches
	go clean -r -i -cache -testcache -modcache

.PHONY: run-client
run-client: ## run client which tries all the api endpoints and prints log output
	CGO_ENABLED=1 go run cmd/client/main.go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
