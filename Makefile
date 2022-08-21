.DEFAULT_GOAL := dev
GOARCH := amd64 # change to arm64 if on mac m1/m2 or surface
PWD := $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

.PHONY: dev
dev: ## generate vet fmt lint test mod-tidy
dev: generate vet fmt lint test mod-tidy

.PHONY: build
build: ## build server binary
build: generate mod-tidy
	CGO_ENABLED=1 GOOS=linux GOARCH=${GOARCH} go build -ldflags='-w -s -extldflags "-static"' -o server cmd/server/main.go

.PHONY: generate
generate: ## go generate and OPENSSL keys generation
	go generate ./...
	openssl req -x509 -nodes -newkey rsa:2048 -keyout resources/certs/server.rsa.key -out resources/certs/server.rsa.crt -days 3650 -subj "/C=CA/ST=QC/O=localhost/CN=localhost"

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
	rm -f coverage.*
	rm -f resources/certs/server.*

.PHONY: run-client
run-client: ## run client which tries all the api endpoints and prints log output
	CGO_ENABLED=1 go run cmd/client/main.go

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
