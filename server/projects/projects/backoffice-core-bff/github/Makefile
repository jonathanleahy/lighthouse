PWD=$(shell pwd -L)
IMAGE_NAME=backoffice-core-bff

.PHONY: all test vendor

all: help

help: ## Display help screen
	@echo "Usage:"
	@echo "	make [COMMAND]"
	@echo "	make help \n"
	@echo "Commands: \n"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: fmt ## Use golintci-lint on your project
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --deadline=65s

clean: fmt ## Remove unused files
	rm -f ./cov.out
	rm -rf bin/

test: fmt test-clean ## Run the tests of the project
	go test -cover -p=1 ./...

test-clean: fmt ## Run the clean cache tests of the project
	go clean -testcache

coverage: fmt ## Run the tests of the project and open the coverage in a Browser
	go test -cover -p=1 -covermode=count -coverprofile=cov.out ./...
	go tool cover -html=cov.out

vendor: fmt ## Copy of all packages needed to support builds and tests in the vendor directory
	go mod vendor

tidy: ## Downloads go dependencies
	go mod tidy

build: ## Build binary
	go build -o application cmd/api/main.go

build-api: fmt ## Build API image
	docker build -t ${IMAGE_NAME}-api -f scripts/docker/api/Dockerfile .

run-docker-api: fmt ## Run API
	docker run --rm -p 8080:8080 --name ${IMAGE_NAME} ${IMAGE_NAME}:latest

run-api: fmt ## Run API
	go run cmd/api/main.go

fmt: tidy ## Run go fmt
	go fmt ./...

openapi-download: ## Download openapi dependencies
	go install github.com/swaggo/swag/cmd/swag@latest

openapi-gen: ## Generate openapi documentation
	go get -u github.com/urfave/cli/v2
	~/go/bin/swag init -g ./internal/app/infrastructure/server/server.go -o ./docs/openapi

graph-gen: ## Run generate graphql resolvers
	go get github.com/99designs/gqlgen@v0.17.49
	go run github.com/99designs/gqlgen generate --verbose

build-mock: ## Build mocks
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen@v1.6.0
	@~/go/bin/mockgen -source=internal/app/infrastructure/network/interfaces.go -destination=internal/app/infrastructure/network/mock/interfaces.go -package=mock
	@~/go/bin/mockgen -source=internal/app/domain/service/interfaces.go -destination=internal/app/domain/service/mock/interfaces.go -package=mock
