#!/usr/bin/env bash

__root_dir=$(dirname "$0")/..

# Fill in revision information
git rev-parse --short HEAD >"$__root_dir"/internal/version/REVISION.txt

# Copy the version info in the root dir
cp "$__root_dir"/VERSION.txt "$__root_dir"/internal/version/VERSION.txt
