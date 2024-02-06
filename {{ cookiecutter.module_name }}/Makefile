.PHONY: bootstrap
bootstrap:
	$(if $(shell command -v mise 2> /dev/null),,$(error ERROR: `mise` [https://mise.jdx.dev] is required))
	mise plugin add https://github.com/olofvndrhr/asdf-just
	mise plugin add https://github.com/asdf-community/asdf-direnv
	mise use just@latest direnv@latest
	go mod tidy
	go install github.com/evilmartians/lefthook@latest
