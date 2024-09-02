// Package context provides service wide shared context values.
// Only add context specific values that are used across all services.
package context

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-workshops/ppp/pkg/logging"
)

type key int

const (
	loggerCtxKey key = iota
)

// WithLogger stores a *zap.Logger inside a given context.
// Use this only when you want to create a sub logger using the With method,
// to populate some fields available across the entire HTTP service.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// Logger retrieves *zap.Logger from a given context.
func Logger(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerCtxKey).(*zap.Logger)
	if ok {
		return logger
	}

	return logging.GetLogger()
}
