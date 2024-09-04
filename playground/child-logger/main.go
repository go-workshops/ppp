package main

import (
	"context"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

func main() {
	ctx := context.Background()
	logger := sharedContext.Logger(ctx).With(zap.String("main_key", "main_value"))
	ctx = sharedContext.WithLogger(ctx, logger)

	logger.Info("logging in main layer")

	l1(ctx)
}

func l1(ctx context.Context) {
	logger := sharedContext.Logger(ctx).With(zap.String("l1_key", "l1_value"))
	ctx = sharedContext.WithLogger(ctx, logger)

	logger.Info("logging in layer 1")
	l2(ctx)
}

func l2(ctx context.Context) {
	logger := sharedContext.Logger(ctx).With(zap.String("l2_key", "l2_value"))
	ctx = sharedContext.WithLogger(ctx, logger)

	logger.Info("logging in layer 2")

	l3(ctx)
}

func l3(ctx context.Context) {
	logger := sharedContext.Logger(ctx).With(zap.String("l3_key", "l3_value"))
	logger.Info("logging in layer 3")
}
