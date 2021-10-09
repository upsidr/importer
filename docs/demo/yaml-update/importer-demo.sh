#!/usr/bin/env bash

demo_helper_type_speed=3000

# shellcheck source=../demo-helper.sh
source "$(dirname "$0")/../demo-helper.sh"

clear_terminal
read -rs
execute "ls -la"

execute "cat demo.yaml"

clear_terminal

execute "cat snippets.yaml"

clear_terminal

comment "Run Importer to update demo.yaml"
execute "importer update demo.yaml"

comment "Update complete, check out the updated file"
execute "cat demo.yaml"
