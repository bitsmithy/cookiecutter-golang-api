#!/usr/bin/env bash

direnv allow
make install
git init && git add . && git commit -m "Initial commit, generated with cookiecutter"
