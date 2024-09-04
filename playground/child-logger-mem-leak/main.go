package main

import (
	"context"
	"time"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

func main() {
	ctx := context.Background()
	logger := sharedContext.Logger(ctx)

	// simulate a streaming service
	for {
		id := time.Now().UnixNano()
		logger = logger.With(zap.Int64("id", id))
		worker(sharedContext.WithLogger(ctx, logger))
	}
}

func worker(ctx context.Context) {
	logger := sharedContext.Logger(ctx)
	logger.Info("worker started")

	// simulate some work
	time.Sleep(1 * time.Second)

	logger.Info("worker finished")
}
