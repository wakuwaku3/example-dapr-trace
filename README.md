# example-dapr-trace

## installation

vscode の dev container で開いてください。

## commands

### 初期化する

```shell
./.cmd/init.sh
```

### server を動かす

```shell
cd server
dapr run --app-port 6006 --app-id server --app-protocol http --dapr-http-port 3501 --runtime-path ../ -- go run .
```

### client を動かす

```shell
cd client
dapr run --app-id client --app-protocol http --dapr-http-port 3500 --runtime-path ../ -- go run .
```

### 全部動かす

```shell
dapr run -f .
```

## Tools

- [zipkin](http://localhost:9411)
- [jaeger](http://localhost:16686)
- [prometheus](http://localhost:9090)
- [grafana](http://localhost:3000)
- [kibana](http://localhost:5601)
