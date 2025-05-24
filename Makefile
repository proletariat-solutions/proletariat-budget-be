APP_NAME=tmpl-hexa-rest-api
TAG=0.0.1
GO_LINT=golangci/golangci-lint

.PHONY: update
## updates: update dependencies
update:
	go get -u ./...

.PHONY: lint
## lint: run linters and static analysis tools to detect potential issues, bugs, and code style violations in Go codebases
lint:
	go mod vendor
	docker run -t --rm -v $(shell pwd):/app -w /app ${GO_LINT} golangci-lint run --fix -v

.PHONY: codegen
## codegen: installs (if necessary) and updates the oapi specs
codegen:
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	oapi-codegen -config ./openapi/oapi-gen.cfg.yaml ./openapi/api-v1.yaml
	go mod tidy

.PHONY: test
## test: run unit tests
test:
	go test -v ./...

.PHONY: run
## run: local with default settings
run:
	docker compose -f examples/docker-compose.yaml up -d && go run main.go

.PHONY: down
## down: stop and remove docker containers
down:
	docker compose -f examples/docker-compose.yaml down

# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
