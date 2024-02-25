.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run:  ### run
	go run ./cmd/app --config-path="./config.yml"
.PHONY: run

update-requirements:  ### update-requirements
	go mod tidy
.PHONY: update-requirements

fmt: ### check go fmt
	gofmt -s -w .
.PHONY: fmt

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

run-service: ### run service
	docker-compose -f docker-compose.yml up --build
.PHONY: run-service
