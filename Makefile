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
	npx @redocly/cli bundle ./openapi/openapi.yaml > ./openapi/out.yaml
	oapi-codegen -config ./openapi/oapi-gen.cfg.yaml ./openapi/out.yaml
	go mod tidy && \
    rm ./openapi/out.yaml

.PHONY: test
## test: run unit integration_tests
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

.PHONY: migrate-create migrate-up migrate-down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations/mysql -seq $$name

migrate-up:
	migrate -path migrations/mysql -database "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}" up

migrate-down:
	migrate -path migrations/mysql -database "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}" down

# help: show this help message
help: Makefile
	@echo
	@echo " Choose a command to run in "$(APP_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
