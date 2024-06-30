package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wakuwaku3/example-dapr-trace/server/lib/errorsx"
	"github.com/wakuwaku3/example-dapr-trace/server/lib/httpx"
	"github.com/wakuwaku3/example-dapr-trace/server/lib/otelx"
	"go.opentelemetry.io/otel/attribute"
)

func getOrder(w http.ResponseWriter, r *http.Request) error {
	logger := otelx.NewLogger(r.Context())
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(errorsx.Wrap(err))
	}

	logger.Info("get order", attribute.KeyValue{
		Key:   "data",
		Value: attribute.StringValue(string(data)),
	})

	logger.Error(errorsx.Wrap(errors.New("example error")))

	_, err = w.Write(data)
	if err != nil {
		logger.Error(errorsx.Wrap(err))
	}

	return nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer cancel()

	shutdown, err := otelx.NewBuilder().
		WithTraceProviderFactory(&otelx.OtlpTraceProviderFactory{
			ServiceName: "server",
			Context:     ctx,
		}).
		WithMeterProviderFactory(&otelx.OtlpMeterProviderFactory{
			ServiceName: "server",
			Context:     ctx,
		}).
		// WithMeterProviderFactory(&otelx.PrometheusMeterProviderFactory{
		// 	ServiceName: "server",
		// }).
		// WithLoggerProviderFactory(&otelx.OtlpLoggerProviderFactory{
		// 	ServiceName: "server",
		// 	Context:     ctx,
		// }).
		Build(ctx)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	server := httpx.NewServer(
		&httpx.ServerOption{
			Port:          "6006",
			ReadTimeout:   5 * time.Second,
			WriteTimeout:  10 * time.Second,
			IdleTimeout:   15 * time.Second,
			CancelTimeout: 5 * time.Second,
		},
	)
	server.HandleFunc("/orders", getOrder)

	if err := server.Serve(ctx); !errors.Is(err, http.ErrServerClosed) {
		log.Println("Error starting HTTP server")
	}
}
