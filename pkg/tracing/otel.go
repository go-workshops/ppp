package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
	"github.com/go-workshops/ppp/pkg/logging"
)

// Supported exporters.
const (
	TraceExporterAttribute = "exporter"

	OTLPTraceExporterName = "otlp"
)

const (
	TraceIDHeader = "trace-id"
	SpanIDHeader  = "span-id"
)

// ServiceName represents the name of the instrumented service.
var ServiceName string

// Tracing errors.
var (
	ErrMissingOTLPURL     = fmt.Errorf("otlp url is required")
	ErrMissingServiceName = fmt.Errorf("service name is required")
)

// TracerProviderConfig represents the distributed tracer provider configuration
type TracerProviderConfig struct {
	TracingEnabled bool
	SpanExporter   SpanExporterWithOptions
	ServiceName    string
	BatchTimeout   time.Duration
	ExportTimeout  time.Duration
	MaxBatchSize   int
	MaxQueueSize   int
}

// SpanExporterWithOptions represents a wrapper around a span exporter with additional resource options per exporter.
type SpanExporterWithOptions struct {
	SpanExporter    traceSDK.SpanExporter
	ResourceOptions []resource.Option
}

// SpanExporter represents a span exporter for Open Telemetry.
type SpanExporter interface {
	name() string
	spanExporter() (SpanExporterWithOptions, error)
}

// NewOTLPExporter represents the OTLP distributed tracing span exporter.
func NewOTLPExporter(url string, timeout ...time.Duration) (SpanExporterWithOptions, error) {
	t := 5 * time.Second
	if len(timeout) > 0 {
		t = timeout[0]
	}

	exporter := otlpTraceExporter{
		url:     url,
		timeout: t,
	}

	return exporter.spanExporter()
}

type otlpTraceExporter struct {
	url     string
	timeout time.Duration
}

func (e otlpTraceExporter) name() string {
	return OTLPTraceExporterName
}

func (e otlpTraceExporter) spanExporter() (SpanExporterWithOptions, error) {
	url := e.url
	if url == "" {
		return SpanExporterWithOptions{}, ErrMissingOTLPURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	conn, err := grpc.NewClient(e.url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return SpanExporterWithOptions{}, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return SpanExporterWithOptions{}, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	se := SpanExporterWithOptions{
		SpanExporter: exporter,
		ResourceOptions: []resource.Option{
			resource.WithAttributes(attribute.String(TraceExporterAttribute, e.name())),
		},
	}

	logging.GetLogger().Info("using the otlp span exporter", zap.String("url", url))
	return se, nil
}

// NewTracerProvider creates a new distributed tracing provider.
// If TracingEnabled is false, it will create a no-op provider.
func NewTracerProvider(cfg TracerProviderConfig) (*Provider, error) {
	if !cfg.TracingEnabled || cfg.SpanExporter.SpanExporter == nil {
		provider := &Provider{
			TracerProvider: noProvider{},
		}
		return provider, nil
	}

	if cfg.ServiceName == "" {
		return nil, ErrMissingServiceName
	}

	ctx := context.Background()
	ServiceName = cfg.ServiceName
	mainResource, err := resource.New(
		ctx,
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
	)
	if err != nil {
		return nil, err
	}
	spanExporterResource, err := resource.New(ctx, cfg.SpanExporter.ResourceOptions...)
	if err != nil {
		return nil, err
	}

	var composedResource *resource.Resource
	for _, r := range []*resource.Resource{
		mainResource,
		spanExporterResource,

		// This resource HAS TO BE THE LAST ONE, otherwise service.name will be "unknown"
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		),
	} {
		composedResource, err = resource.Merge(composedResource, r)
		if err != nil {
			return nil, err
		}
	}

	tracerProvider := traceSDK.NewTracerProvider(
		traceSDK.WithBatcher(
			cfg.SpanExporter.SpanExporter,
			traceSDK.WithBatchTimeout(cfg.BatchTimeout),
			traceSDK.WithExportTimeout(cfg.ExportTimeout),
			traceSDK.WithMaxExportBatchSize(cfg.MaxBatchSize),
			traceSDK.WithMaxQueueSize(cfg.MaxQueueSize),
		),

		// This MUST be a COMPOSED resource, otherwise only the LAST CALL to WithResource() will be considered,
		// even if this is a variadic function :P.
		traceSDK.WithResource(composedResource),
	)
	provider := &Provider{
		TracerProvider: tracerProvider,
	}

	return provider, nil
}

// Provider represents a wrapper around traceSDK.TracerProvider
// which has more methods such as Shutdown. Unfortunately the
// trace.TracerProvider does not have a Shutdown method.
type Provider struct {
	TracerProvider
}

// ForceFlush is a wrapper around traceSDK.TracerProvider.ForceFlush
func (p *Provider) ForceFlush(ctx context.Context) error {
	return p.TracerProvider.ForceFlush(ctx)
}

// Shutdown is a wrapper around traceSDK.TracerProvider.Shutdown
func (p *Provider) Shutdown(ctx context.Context) error {
	return p.TracerProvider.Shutdown(ctx)
}

// TracerProvider represents both Provider and traceSDK.TracerProvider.
// We need this to be able to call Shutdown on application shutdown.
// The trace.TracerProvider interface does not contain the Shutdown method.
type TracerProvider interface {
	Tracer(string, ...trace.TracerOption) trace.Tracer
	ForceFlush(ctx context.Context) error
	Shutdown(ctx context.Context) error
	embedded.TracerProvider
}

type noProvider struct {
	embedded.TracerProvider
}

func (noProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	return noTracer{}
}

func (noProvider) ForceFlush(context.Context) error {
	return nil
}

func (noProvider) Shutdown(context.Context) error {
	return nil
}

type noTracer struct {
	embedded.Tracer
}

func (noTracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, noSpan{}
}

type noSpan struct {
	embedded.Span
}

func (s noSpan) AddLink(trace.Link) {
}

func (noSpan) End(...trace.SpanEndOption) {
}

func (noSpan) AddEvent(string, ...trace.EventOption) {
}

func (noSpan) IsRecording() bool {
	return false
}

func (noSpan) RecordError(error, ...trace.EventOption) {
}

func (noSpan) SpanContext() trace.SpanContext {
	return trace.SpanContext{}
}

func (noSpan) SetStatus(codes.Code, string) {
}

func (noSpan) SetName(string) {
}

func (noSpan) SetAttributes(...attribute.KeyValue) {
}

func (noSpan) TracerProvider() trace.TracerProvider {
	return noProvider{}
}

// NewTextMapPropagator represents a custom HTTP Inject/Extract propagator, used
// to inject/extract the trace-id and span-id into/from the HTTP Headers.
func NewTextMapPropagator(ctx context.Context) propagation.TextMapPropagator {
	logger := sharedContext.Logger(ctx).With(zap.String("source", "open_telemetry"))
	return propagator{logger: logger}
}

type propagator struct {
	logger *zap.Logger
	propagation.TraceContext
}

func (p propagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return
	}

	// Pick one :P
	// Custom headers
	carrier.Set(TraceIDHeader, sc.TraceID().String())
	carrier.Set(SpanIDHeader, sc.SpanID().String())
	// Default headers
	p.TraceContext.Inject(ctx, carrier)
}

func (p propagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	// Pick one :P
	// Custom headers
	traceIDHeader := carrier.Get(TraceIDHeader)
	spanIDHeader := carrier.Get(SpanIDHeader)
	// Default headers
	traceContextCtx := p.TraceContext.Extract(ctx, carrier)
	return sharedContext.WithSpanContext(traceContextCtx, traceIDHeader, spanIDHeader)
}

func (p propagator) Fields() []string {
	return append(
		// Pick one :P
		// Default headers
		p.TraceContext.Fields(),
		// Custom headers
		TraceIDHeader, SpanIDHeader,
	)
}
