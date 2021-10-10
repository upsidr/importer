#!/usr/bin/env bash

demo_helper_type_speed=3000

# shellcheck source=../demo-helper.sh
source "$(dirname "$0")/../demo-helper.sh"

clear_terminal
read -rs
execute "ls -la"

execute "cat demo.yaml"
execute "cat snippets.yaml"

clear_terminal

comment "Run Importer Purge on 'demo.yaml'"
execute "importer purge demo.yaml"

comment "Confirm lines surrounded by Importer Markers are now purged."
execute "cat demo.yaml"
