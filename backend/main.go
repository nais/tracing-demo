package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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

type Config struct {
	ListenAddress string `envconfig:"LISTEN_ADDR" default:"127.0.0.1:8080"`
	TraceName     string `envconfig:"TRACE_NAME" default:"github.com/nais/tracing-demo/backend"`
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

	tp, err := newProvider(os.Stdout)
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

	err = http.ListenAndServe(cfg.ListenAddress, otelHandler)
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	prop := otel.GetTextMapPropagator()
	ctx := prop.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	//prop.Inject(r.Context(), propagation.HeaderCarrier(w.Header()))

	parentSpan := trace.SpanFromContext(ctx)
	link := trace.LinkFromContext(ctx)
	ctx = trace.ContextWithSpan(ctx, parentSpan)
	//parentSpan.AddEvent("foobar happened")

	tracer := otel.Tracer(h.Config.TraceName)
	ctx, span := tracer.Start(ctx, "ServeHTTP", trace.WithLinks(link))
	defer span.End()
	span.AddEvent("world burned")

	fmt.Printf("my header was: %s\n", r.Header.Get("Traceparent"))
	span.RecordError(fmt.Errorf("Parent traceid=%v spanid=%v", link.SpanContext.TraceID(), link.SpanContext.SpanID()))

	request := &Request{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	num, err := Fibonacci(ctx, request.Number)
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

// Fibonacci returns the n-th fibonacci number.
func Fibonacci(ctx context.Context, n uint) (uint, error) {
	if n <= 1 {
		return uint(n), nil
	}

	var n2, n1 uint = 0, 1
	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}

func newProvider(w io.Writer) (*sdk_trace.TracerProvider, error) {
	exp, err := newExporter(w)
	if err != nil {
		return nil, err
	}

	return sdk_trace.NewTracerProvider(
		sdk_trace.WithBatcher(exp),
		sdk_trace.WithResource(newResource()),
	), nil
}

// newExporter returns a console exporter.
func newExporter(w io.Writer) (sdk_trace.SpanExporter, error) {
	return otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithEndpoint("localhost:4317"),
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
