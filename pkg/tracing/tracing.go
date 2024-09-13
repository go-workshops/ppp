// Package tracing provides tracing utilities and Open Telemetry setup across all Platform services.
package tracing

import (
	"context"
	"runtime/debug"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

// Open Telemetry tracing attribute keys
const (
	DBAttributeKey        = "db"
	StreamingAttributeKey = "streaming"
	HTTPAttributeKey      = "http"
)

// Open Telemetry tracing attribute values
const (
	PostgresAttributeValue = "postgresql"
	RedisAttributeValue    = "redis"
	KafkaAttributeValue    = "kafka"
	K8SAttributeValue      = "k8s"
)

// StartPostgres creates a new Postgres Open Telemetry tracing Span.
// Use this every time you want to trace a Postgres database operation.
func StartPostgres(ctx context.Context, traceName string) (context.Context, trace.Span) {
	tr := otel.Tracer(traceName)
	attributes := withServiceName(
		ctx,
		operationNameAttribute(traceName),
		attribute.String(DBAttributeKey, PostgresAttributeValue),
	)

	ctx, span := tr.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attributes...)

	return ctx, span
}

// StartRedis creates a new Redis Open Telemetry tracing Span.
// Use this every time you want to trace a Redis database operation.
func StartRedis(ctx context.Context, traceName string) (context.Context, trace.Span) {
	tr := otel.Tracer(traceName)
	attributes := withServiceName(
		ctx,
		operationNameAttribute(traceName),
		attribute.String(DBAttributeKey, RedisAttributeValue),
	)

	ctx, span := tr.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attributes...)

	return ctx, span
}

// StartKafkaProducer creates a new Kafka Producer Open Telemetry tracing Span.
// Use this every time you want to trace a Kafka producer operation.
func StartKafkaProducer(ctx context.Context, traceName string) (context.Context, trace.Span) {
	tr := otel.Tracer(traceName)
	attributes := withServiceName(
		ctx,
		operationNameAttribute(traceName),
		attribute.String(StreamingAttributeKey, KafkaAttributeValue),
	)

	ctx, span := tr.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindProducer))
	span.SetAttributes(attributes...)

	return ctx, span
}

// StartKafkaConsumer creates a new Kafka Consumer Open Telemetry tracing Span.
// Use this every time you want to trace a Kafka consumer operation.
func StartKafkaConsumer(ctx context.Context, traceName string) (context.Context, trace.Span) {
	tr := otel.Tracer(traceName)
	attributes := withServiceName(
		ctx,
		operationNameAttribute(traceName),
		attribute.String(StreamingAttributeKey, KafkaAttributeValue),
	)

	ctx, span := tr.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindConsumer))
	span.SetAttributes(attributes...)

	return ctx, span
}

// StartHTTP creates a new HTTP Open Telemetry tracing Span.
// Use this every time you want to trace an external HTTP call operation.
func StartHTTP(ctx context.Context, httpServiceName, traceName string) (context.Context, trace.Span) {
	tr := otel.Tracer(traceName)
	attributes := withServiceName(
		ctx,
		operationNameAttribute(traceName),
		attribute.String(HTTPAttributeKey, httpServiceName),
	)

	ctx, span := tr.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attributes...)

	return ctx, span
}

// StartK8S creates a new K8S Open Telemetry tracing Span.
// Use this every time you want to trace a K8S SKD operation.
func StartK8S(ctx context.Context, traceName string) (context.Context, trace.Span) {
	return StartHTTP(ctx, K8SAttributeValue, traceName)
}

// RecordError records an error in the Open Telemetry tracing Span also adding the stacktrace to the recorded error.
func RecordError(ctx context.Context, err error, description string) {
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, description)

	if err != nil {
		span.RecordError(err, trace.WithAttributes(semconv.ExceptionStacktraceKey.String(string(debug.Stack()))))
	}
}

func operationNameAttribute(operationName string) attribute.KeyValue {
	return attribute.String("operation_name", operationName)
}

func withServiceName(ctx context.Context, attributes ...attribute.KeyValue) []attribute.KeyValue {
	attrs := append([]attribute.KeyValue{}, attributes...)
	serviceName := sharedContext.ServiceName(ctx)
	if serviceName != "" {
		attrs = append(attrs, semconv.ServiceNameKey.String(serviceName))
	}

	return attrs
}
