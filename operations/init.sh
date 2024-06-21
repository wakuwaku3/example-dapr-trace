script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")

dapr uninstall
dapr init

$script_dir/build.sh
