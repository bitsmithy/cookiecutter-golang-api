.PHONY: bootstrap
bootstrap:
	$(if $(shell command -v mise 2> /dev/null),,$(error ERROR: `mise` [https://mise.jdx.dev] is required))
	mise plugin add just
	mise plugin add direnv
	mise install just@latest direnv@latest
	mise use just@latest direnv@latest
	mise exec golang -- go mod tidy
	mise exec golang -- go install github.com/evilmartians/lefthook@latest
	mise exec golang -- go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	mise exec golang -- go install gotest.tools/gotestsum@latest
	mise exec golang -- go install github.com/cosmtrek/air@latest
