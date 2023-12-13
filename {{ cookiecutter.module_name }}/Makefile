include .envrc

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.envrc:
	@echo -e "{{ cookiecutter.module_name.upper() }}_ENV=local\n{{ cookiecutter.module_name.upper() }}_PORT={{ cookiecutter.server_port }}" > $@

## install: bootstrap all dependencies for the project
.PHONY: install
install: .envrc
	@asdf --version >/dev/null 2>&1 || (echo "ERROR: asdf (https://asdf-vm.com) is required. Install and try again."; exit 1)
	asdf install
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install gotest.tools/gotestsum@latest

## run: start the API server
.PHONY: run
run: .envrc
	{{ cookiecutter.module_name.upper() }}_ENV={{ "$" }}{{ "{" }}{{ cookiecutter.module_name.upper() }}_ENV{{ "}" }} \
	go run ./cmd/api

## dependencies: install all code dependencies
.PHONY: dependencies
dependencies: .envrc
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify

## test: run all tests
.PHONY: test
test: .envrc dependencies
	@echo 'Running tests...'
	go clean -testcache
	gotestsum --format dots-v2 -- -race -vet=off ./...

## audit: tidy dependencies and format, and vet all code
.PHONY: audit
audit: .envrc dependencies
	@echo 'Auditing code...'
	golangci-lint run

## deploy: deploy the API
.PHONY: deploy
deploy: audit test confirm
	@echo 'Nothing to do: implement `make deploy`'
