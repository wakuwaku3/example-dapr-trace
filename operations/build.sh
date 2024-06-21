script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")

cd $script_dir/../client
go build -o ./out/ .

cd $script_dir/../server
go build -o ./out/ .

printf "Build complete\n"
