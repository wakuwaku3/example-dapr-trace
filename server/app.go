package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wakuwaku3/example-dapr-trace/lib/fluentx"
	"github.com/wakuwaku3/example-dapr-trace/lib/httpx"
	"github.com/wakuwaku3/example-dapr-trace/lib/logx"
	"github.com/wakuwaku3/example-dapr-trace/lib/otelx"
	"github.com/wakuwaku3/example-dapr-trace/server/app/order"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer cancel()

	shutdown, err := otelx.NewBuilder("dapr-diagnostics").
		WithTraceProviderFactory(&otelx.OtlpTraceProviderFactory{
			ServiceName: "server",
			Context:     ctx,
		}).
		WithMeterProviderFactory(&otelx.OtlpMeterProviderFactory{
			ServiceName: "server",
			Context:     ctx,
		}).
		WithLoggerProviderFactory(&otelx.StdoutLoggerProviderFactory{
			ServiceName: "server",
			Context:     ctx,
		}).
		Build(ctx)

	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	logx.Provider.Set(otelx.NewLogger())
	fluentLogger, err := fluentx.NewLogger(fluentx.Config{
		FluentHost: "fluentd",
		FluentPort: 24224,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := fluentLogger.Close(); err != nil {
			log.Fatalf("failed to close fluent logger: %v", err)
		}
	}()
	logx.Provider.Set(fluentLogger)

	server := httpx.NewServer(
		&httpx.ServerOption{
			Port:          "6006",
			ReadTimeout:   5 * time.Second,
			WriteTimeout:  10 * time.Second,
			IdleTimeout:   15 * time.Second,
			CancelTimeout: 5 * time.Second,
		},
	)

	// route.
	server.HandleFunc("/orders", order.Get)

	if err := server.Serve(ctx); !errors.Is(err, http.ErrServerClosed) {
		log.Println("Error starting HTTP server")
	}
}
