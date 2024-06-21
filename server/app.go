package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Llongfile)

func getOrder(w http.ResponseWriter, r *http.Request) {
	tr := otel.GetTracerProvider().Tracer("dapr-diagnostics")
	log.Println(r.Header)

	traceparent := r.Header.Get("Traceparent")
	tid, err := trace.TraceIDFromHex(strings.Split(traceparent, "-")[1])
	if err != nil {
		log.Println("Error parsing traceparent:", err.Error())
	}
	sid, err := trace.SpanIDFromHex(strings.Split(traceparent, "-")[2])
	if err != nil {
		log.Println("Error parsing traceparent:", err.Error())
	}

	ctx, span := tr.Start(trace.ContextWithSpanContext(r.Context(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: trace.FlagsSampled,
	}).WithTraceID(tid)), r.Method+" "+r.RequestURI, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(semconv.HTTPMethodKey.String(r.Method), semconv.HTTPURLKey.String(r.RequestURI)))

	r = r.WithContext(ctx)
	log.Println(ctx)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err.Error())
	}
	defer span.End()

	_, span2 := tr.Start(ctx, r.Method+" "+r.RequestURI, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(semconv.HTTPMethodKey.String(r.Method), semconv.HTTPURLKey.String(r.RequestURI)))
	span2.End()

	span.AddEvent(fmt.Sprintln("Order received:", string(data)))

	logger.Println("Order received:", string(data))
	_, err = w.Write(data)
	if err != nil {
		log.Println("Error writing the response:", err.Error())
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	url := "http://localhost:9411/api/v2/spans"
	shutdown, err := initTracer(url, logger)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: %w", err)
		}
	}()

	// Create a new router and respond to POST /orders requests
	r := mux.NewRouter()
	r.HandleFunc("/orders", getOrder).Methods("POST")

	// Start the server listening on port 6006
	// This is a blocking call
	if err := http.ListenAndServe(":6006", r); !errors.Is(err, http.ErrServerClosed) {
		log.Println("Error starting HTTP server")
	}
}

func initTracer(url string, logger *log.Logger) (func(context.Context) error, error) {
	// Create Zipkin Exporter and install it as a global tracer.
	//
	// For demoing purposes, always sample. In a production application, you should
	// configure the sampler to a trace.ParentBased(trace.TraceIDRatioBased) set at the desired
	// ratio.
	exporter, err := zipkin.New(
		url,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("server"),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.Shutdown, nil
}
