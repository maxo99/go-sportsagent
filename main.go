package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sportsagent/internal/handlers"
	"sportsagent/internal/version"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func setupServer() *http.ServeMux {
	mux := http.NewServeMux()
	handler := handlers.NewAgentHandler()
	toolsHandler := handlers.NewToolsHandler()
	mux.Handle("/query", otelhttp.NewHandler(http.HandlerFunc(handler.HandleQuery), "Query"))
	mux.Handle("/tools", otelhttp.NewHandler(http.HandlerFunc(toolsHandler.HandleGetTools), "Tools"))
	mux.HandleFunc("/healthz", handlers.HandleHealth)
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}

func main() {
	godotenv.Load()

	shutdown, err := initTracing(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if err != nil {
		log.Fatalf("setup tracing: %v", err)
	}
	defer shutdown(context.Background())

	mux := setupServer()

	log.Println("Starting GoSportsAgent version:", version.Version, "server on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}

func initTracing(endpoint string) (func(context.Context) error, error) {
	exp, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpointURL(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-sportsagent"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.Shutdown, nil
}
