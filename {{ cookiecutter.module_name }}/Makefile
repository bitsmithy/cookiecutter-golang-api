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
	asdf install
	go install honnef.co/go/tools/cmd/staticcheck@latest
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
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	golangci-lint run

## deploy: deploy the API
.PHONY: deploy
deploy: audit test confirm
	@echo 'Nothing to do: implement `make deploy`'
