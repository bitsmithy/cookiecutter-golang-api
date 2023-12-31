#!/usr/bin/env bash

direnv allow
asdf install golang latest
asdf local golang latest
go mod init '{{ cookiecutter.module_path }}'
make setup
make audit
git init && git add . && git commit -m "Initial commit, generated with cookiecutter"
lefthook install
