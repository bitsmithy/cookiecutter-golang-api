include .envrc

BUILD_OUTPUT := /tmp/bin/{{ cookiecutter.module_name }}

## help: print this help message (keep this as the first target to make it default)
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: check_all
check_all: format check_dirty_workspace lint audit test ## check_dirty_workspace must come after format to exit on error

.PHONY: check_dirty_workspace
check_dirty_workspace:
	@git diff --quiet

.PHONY: confirm
confirm:
	@printf 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## audit: run quality control checks
.PHONY: audit
audit:
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## build: build the app
.PHONY: build
build: check_all
	go build -o $(BUILD_OUTPUT) ./cmd/api

## deps: install all code dependencies
.PHONY: deps
deps:
	go mod tidy -v
	go mod verify

## deploy: deploy the API
.PHONY: deploy
deploy: check_all confirm
	@echo 'No deploy actions yet: implement `make deploy` section that actually deploys'

## format: autoformat all code
.PHONY: format
format:
	go run github.com/incu6us/goimports-reviser/v3@latest -rm-unused -set-alias -format -set-exit-status ./...
	go run mvdan.cc/gofumpt@latest -l -w .

## lint: run all linters
.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --fix

## run: start the API server
.PHONY: run
run: build
	$(BUILD_OUTPUT)

## run/live: run the API server with reloading on file changes
.PHONY: run/live
run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "go build -o $(BUILD_OUTPUT) ./cmd/api" --build.bin "$(BUILD_OUTPUT)" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

## setup: bootstrap all development dependencies for the project
.PHONY: setup
setup:
	$(if $(shell command -v asdf 2> /dev/null),,$(error ERROR: `asdf` [https://asdf-vm.com] is required))
	asdf install
	go mod tidy
	go install github.com/evilmartians/lefthook@latest

## test: run all tests
.PHONY: test
test:  deps
	go clean -testcache
	go run gotest.tools/gotestsum@latest -- -race -buildvcs -vet=off ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go run gotest.tools/gotestsum@latest -- -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out
