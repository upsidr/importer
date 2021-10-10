#!/usr/bin/env bash

demo_helper_type_speed=3000

# shellcheck source=../demo-helper.sh
source "$(dirname "$0")/../demo-helper.sh"

clear_terminal
read -rs
execute "ls -la"

execute "cat demo-template.yaml"
comment "Notice how there is a special marker '== importer-skip-update ==' at the top of the file."

clear_terminal

execute "cat snippets.yaml"

clear_terminal

comment_r "Run Importer Update - but this should have no effect."
execute "importer update demo-template.yaml"

comment_r "Update succeeds, but nothing happens to the file."
execute "cat demo-template.yaml"

clear_terminal

comment "Now, let's use 'demo-template.yaml' for generating a new file."
execute "importer generate demo-template.yaml --out demo.gen.yaml"

comment "You can see that a new file is created."
execute "ls -l"

clear_terminal

comment "Let's check that 'demo-template.yaml' is left intact."
execute "cat demo-template.yaml"

comment "And see what 'demo.gen.yaml' has now."
execute "cat demo.gen.yaml"

comment_g "This time, you can see that the generated file no longer has Importer Markers."
