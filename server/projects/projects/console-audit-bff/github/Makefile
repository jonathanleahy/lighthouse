PWD = $(shell pwd -L)
GOCMD=go
DOCKERCMD=docker
DOCKERCOMPOSECMD=docker-compose
GOTEST=$(GOCMD) test
IMAGE_NAME = console-audit-bff
LIBRARY_ENV ?= dev

.PHONY: all test vendor

all: help

help: ## Display help screen
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help \n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: lint-go lint-yaml ## Run all available linters

lint-go: fmt ## Use golintci-lint on your project
	$(DOCKERCMD) run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s

lint-yaml: fmt ## Use yamllint on the yaml file of your projects
	$(DOCKERCMD) run --rm -it -v $(PWD):/data cytopia/yamllint -f parsable $(shell git ls-files '*.yml' '*.yaml')

clean: fmt ## Remove unused files
	rm -f ./coverage.out
	rm -rf bin/

test: fmt test-clean ## Run the tests of the project
	$(GOTEST) -cover -p=1 ./...

test-clean: fmt ## Run the clean cache tests of the project
	$(GOCMD) clean -testcache

coverage: fmt ## Run the tests of the project and open the coverage in a Browser
	$(GOTEST) -cover -p=1 -covermode=count -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

vendor: fmt ## Copy of all packages needed to support builds and tests in the vendor directory
	$(GOCMD) mod vendor

tidy: ## Downloads go dependencies
	$(GOCMD) mod tidy

docker-compose-up: ## Run docker-compose services of project
	$(DOCKERCOMPOSECMD) up -d

docker-compose-down: ## Stop docker-compose services of project
	$(DOCKERCOMPOSECMD) down --remove-orphans

docker-compose-restart: docker-compose-down docker-compose-up ## Restart docker-compose services of project

docker-compose-logs: ## Logs docker-compose containers of project
	$(DOCKERCOMPOSECMD) logs -f app

docker-compose-ps: ## List docker-compose containers of project
	$(DOCKERCOMPOSECMD) ps

build: build-api ## Build all image

build-api: fmt ## Build API image
	$(DOCKERCMD) build -t ${IMAGE_NAME}-api -f scripts/docker/api/Dockerfile .

run-docker-api: fmt ## Run API
	$(DOCKERCMD) run --rm -p 8080:8080 --name console-bff console-audit-bff-api:latest

run-api: fmt ## Run API
	$(GOCMD) run cmd/api/main.go

fmt: tidy ## Run go fmt
	go fmt ./...

openapi-download: ## Download openapi dependencies
	go get -d github.com/swaggo/swag/cmd/swag

openapi-gen: ## Generate openapi documentation
	go run github.com/swaggo/swag/cmd/swag init -g ./internal/app/adapter/web/router.go -o ./docs/openapi

graph-gen: ## Run generate graphql resolvers
	go run github.com/99designs/gqlgen generate --verbose

graph-download: ## Download gqlgen binary
	go mod tidy
	go get -d github.com/99designs/gqlgen@v0.17.4

download-dependencies: openapi-download graph-download ## Downloads all dependencies

generate-all: openapi-gen graph-gen ## Run all generate

build-mock: ## Build mocks
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen@v1.6.0
	@~/go/bin/mockgen -source=internal/app/domain/service/audit.go -destination=internal/app/domain/service/mock/audit.go -package=mock
	@~/go/bin/mockgen -source=internal/app/domain/service/net.go -destination=internal/app/domain/service/mock/net.go -package=mock
	# @~/go/bin/mockgen -source=internal/app/infrastructure/network/audit.go -destination=internal/app/infrastructure/network/mock/audit.go -package=mock

