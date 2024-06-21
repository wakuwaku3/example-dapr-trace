module github.com/wakuwaku3/example-dapr-trace/server

go 1.22

require (
	go.opentelemetry.io/otel v1.27.0
	go.opentelemetry.io/otel/sdk v1.27.0
)

require (
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/exporters/zipkin v1.27.0
	go.opentelemetry.io/otel/metric v1.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.27.0
)
