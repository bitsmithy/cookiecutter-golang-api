#!/usr/bin/env bash

if ! command -v mise &>/dev/null; then
  echo "ERROR: \`mise\` [https://mise.jdx.dev] is required. Please install it then try again."
  exit 1
fi
