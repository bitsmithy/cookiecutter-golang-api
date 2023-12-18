#!/usr/bin/env bash

if ! command -v asdf &>/dev/null; then
  echo "ERROR: \`asdf\` [https://asdf-vm.com] is required. Please install it then try again."
  exit 1
fi
