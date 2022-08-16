.DEFAULT_GOAL := dev

.PHONY: dev
dev: ## dev build
dev: clean install generate vet fmt lint test mod-tidy

.PHONY: ci
ci: ## CI build
ci: dev diff

.PHONY: clean
clean: ## remove files created during build pipeline
	rm -rf dist
	rm -f coverage.*

.PHONY: install
install: ## go install tools
	cd tools && go install $(shell cd tools && go list -f '{{ join .Imports " " }}' -tags=tools)

.PHONY: generate
generate: ## go generate
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
	cd tools && go mod tidy

.PHONY: diff
diff: ## git diff
	git diff --exit-code
	RES=$$(git status --porcelain) ; if [ -n "$$RES" ]; then echo $$RES && exit 1 ; fi

.PHONY: run
run: ## go run
	@go run -race .

.PHONY: go-clean
go-clean: ## go clean build, test and modules caches
	go clean -r -i -cache -testcache -modcache

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
