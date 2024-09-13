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
)

const (
	DBAttributeKey   = "db"
	HTTPAttributeKey = "http"
)

const (
	PostgresAttributeValue = "postgresql"
)

// StartPostgres creates a new Postgres Open Telemetry tracing Span.
// Use this every time you want to trace a Postgres database operation.
func StartPostgres(ctx context.Context, traceName string) (context.Context, trace.Span) {
	tr := otel.Tracer(traceName)
	attributes := []attribute.KeyValue{
		operationNameAttribute(traceName),
		attribute.String(DBAttributeKey, PostgresAttributeValue),
	}

	ctx, span := tr.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attributes...)

	return ctx, span
}

// StartHTTP creates a new HTTP Open Telemetry tracing Span.
// Use this every time you want to trace an external HTTP call operation.
func StartHTTP(ctx context.Context, httpServiceName, traceName string) (context.Context, trace.Span) {
	tr := otel.Tracer(traceName)
	attributes := []attribute.KeyValue{
		operationNameAttribute(traceName),
		attribute.String(HTTPAttributeKey, httpServiceName),
	}

	ctx, span := tr.Start(ctx, traceName, trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(attributes...)

	return ctx, span
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
