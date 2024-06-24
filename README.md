# example-dapr-trace

## installation

vscode の dev container で開いてください。

```shell
./operations/init.sh
```

## commands

### server を動かす

```shell
cd server
dapr run --app-port 6006 --app-id server --app-protocol http --dapr-http-port 3501 -- go run .
```

### client を動かす

```shell
cd client
dapr run --app-id client --app-protocol http --dapr-http-port 3500 -- go run .
```

### 全部動かす

```shell
dapr run -f .
```

### zipkin にアクセスする

http://localhost:9411 にアクセスする
