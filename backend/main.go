package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Simple HTTP service that does some "work" and attaches it to a trace.

const progName = "tracing-demo-backend"
const tracerName = ""

type Config struct {
	ListenAddress string `envconfig:"LISTEN_ADDR" default:"127.0.0.1:8080"`
	Endpoint      string `envconfig:"ENDPOINT" default:"localhost:4317"`
}

type Request struct {
	Number uint
}

type Response struct {
	Number uint
}

type Handler struct {
	Config Config
}

func main() {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err)
	}

	log.Println("Tracing-demo backend starting up")

	tp, err := newProvider(cfg.Endpoint)
	if err != nil {
		panic(err)
	}

	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	handler := &Handler{
		Config: *cfg,
	}
	otelHandler := otelhttp.NewHandler(handler, "")

	log.Println("Ready to receive requests")
	err = http.ListenAndServe(cfg.ListenAddress, otelHandler)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// This is supposed to look like a real API endpoint in a production app.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Initialize OpenTelemetry using information from the `Traceheader` header field.
	// This way we can visualize API calls across multiple applications.
	prop := otel.GetTextMapPropagator()
	ctx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	link := trace.LinkFromContext(ctx)
	tracer := otel.Tracer(tracerName)

	// Create a new trace and give it a meaningful name.
	ctx, span := tracer.Start(ctx, "Backend API", trace.WithLinks(link))
	defer span.End()

	msg := fmt.Sprintf("Request received with trace header '%s'", r.Header.Get("Traceparent"))
	log.Println(msg)
	span.AddEvent(msg)

	request := &Request{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var num uint
	RunTrace(ctx, "Fibonacci sequence generator", func() {
		num, err = Fibonacci(request.Number)
	})

	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := &Response{
		Number: num,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		span.RecordError(err)
	}
}

func RunTrace(ctx context.Context, name string, fn func()) {
	tracer := otel.Tracer(tracerName)
	_, span := tracer.Start(ctx, name, trace.WithLinks(trace.LinkFromContext(ctx)))
	fn()
	span.End()
}

// Fibonacci returns the n-th fibonacci number.
func Fibonacci(n uint) (uint, error) {
	if n <= 1 {
		return n, nil
	}

	var n2, n1 uint = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}

func newProvider(endpoint string) (*sdk_trace.TracerProvider, error) {
	exp, err := newExporter(endpoint)
	if err != nil {
		return nil, err
	}

	return sdk_trace.NewTracerProvider(
		sdk_trace.WithBatcher(exp),
		sdk_trace.WithResource(newResource()),
	), nil
}

// newExporter returns a console exporter.
func newExporter(endpoint string) (sdk_trace.SpanExporter, error) {
	return otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(progName),
			semconv.ServiceVersion("v1.0.0"),
			attribute.String("environment", "dev"),
		),
	)
	return r
}
