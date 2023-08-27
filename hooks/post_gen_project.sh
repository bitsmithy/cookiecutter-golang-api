#!/usr/bin/env bash

direnv allow
go mod tidy
make install
git init && git add . && git commit -m "Initial commit, generated with cookiecutter"
