script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")

dapr uninstall --runtime-path $script_dir/../
dapr init --runtime-path $script_dir/../

$script_dir/build.sh

docker ps
