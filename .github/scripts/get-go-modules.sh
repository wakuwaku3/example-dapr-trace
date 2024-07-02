#!/bin/bash
script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")

cd $script_dir
root_path=$(builtin cd $script_dir/../..; pwd)
cd $root_path
black_list=".cmd/|.dapr/|.devcontainer/|.github/|.vscode/"
results+=""

for dir in $(ls -d */ | grep -v -E "$black_list"); do
    dir=${dir%/}  # Remove trailing slash
    target_path="$root_path/$dir"
    cmd="$target_path/go.mod"
    if [[ ! -e $cmd ]]; then
        continue
    fi

    if [[ -n $results ]]; then
        results+="\n"
    fi
    results+="$dir: "
    results+="\n  - '$dir/**/*'"
done

echo -e $results
