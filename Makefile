.PHONY: install-migrate
## install-migrate: installs the golang migrate binary
install-migrate:
	./scripts/install_migrate.sh

.PHONY: build
## build: Builds the Go program
build:
	CGO_ENABLED=0 go build -o seeder ./cli

.PHONY: up
## up: Runs all the containers listed in the docker-compose.yml file
up: build
	docker-compose up --build -d

.PHONY: down
## down: Shut down all the containers listed in the docker-compose.yml file
down:
	docker-compose down

.PHONY: linter
## linter: Runs the golangci-lint command
linter:
	golangci-lint run --enable=golint --enable=godot ./...

.PHONY: test
## test: Runs all the test suite of the project
test:
	go test -v -race ./...

.PHONY: unit-test
## unit-test: Runs the tests with the short flag
unit-test:
	go test -v -short -race ./...

.PHONY: int-test
## int-test: Runs the integration tests in a docker environment
int-test:
	docker-compose run --entrypoint=make seeder test

.PHONY: migrate
## migrate: Runs the migrations
migrate:
	migrate -source file://db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

.PHONY: run
## run: Runs the application
run: build
	./seeder

.PHONY: help
## help: Prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
