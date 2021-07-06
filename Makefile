.PHONY: install-migrate
## install-migrate: installs the golang migrate binary
install-migrate:
	./scripts/install_migrate.sh

.PHONY: build
## build: Builds the Go program
build:
	go build -o seeder .

run: build
.PHONY: up
## up: Runs all the containers listed in the docker-compose.yml file
up:
	docker-compose up --build -d

.PHONY: down
## down: Shut down all the containers listed in the docker-compose.yml file
down:
	docker-compose down
