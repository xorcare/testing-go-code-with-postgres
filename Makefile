.DEFAULT_GOAL := help

COVER_FILE ?= coverage.out

.PHONY: build
build: ## Build a command to quickly check compiles.
	@go build ./...

.PHONY: check
check: lint build test ## Runs all necessary code checks.

.PHONY: test
test: ## Run all tests.
	@go test -race -count=1 -coverprofile=$(COVER_FILE) ./...
	@go tool cover -func=$(COVER_FILE) | grep ^total | tr -s '\t'

.PHONY: test-short
test-short: ## Run only unit tests, tests without I/O dependencies.
	@go test -short ./...

.PHONY: test-env-up
test-env-up: ## Run test environment.
	@docker-compose up migrate

.PHONY: test-env-down
test-env-down: ## Down and cleanup test environment.
	@docker-compose down -v

.PHONY: lint
lint: tools ## Check the project with lint.
	@golangci-lint run --fix ./...

tools: ## Install all needed tools, e.g.
	@go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.2

.PHONY: help
help: ## Show help for each of the Makefile targets.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
