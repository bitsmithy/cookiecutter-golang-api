#!/usr/bin/env bash

direnv allow
mise install golang@latest
mise use golang@latest
go mod init '{{ cookiecutter.module_path }}'
make bootstrap
just lint
just test
git init && git add . && git commit -m "Initial commit, generated with cookiecutter"
lefthook install
