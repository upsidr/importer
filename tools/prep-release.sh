#!/usr/bin/env bash

__root_dir=$(dirname "$0")/..

# Fill in revision information
git rev-parse --short HEAD >"$__root_dir"/internal/version/REVISION.txt

# Assuming that tag value is provided as an argument, use that as the version
# information. If missing by any chance, this defaults to "unknown".
__tag=$1
if [[ -z $__tag ]]; then
    __tag="unknown"
fi
echo "$__tag" >"$__root_dir"/internal/version/VERSION.txt
