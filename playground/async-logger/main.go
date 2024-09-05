package main

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.OutputPaths = []string{"zap.log"}

	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return l
}

func req(wg *sync.WaitGroup, id int, logger *zap.Logger, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }()
	logger.Info("request", zap.Int("req_id", id))
}

func main() {
	l := logger()
	defer func() { _ = l.Sync() }()

	var wg sync.WaitGroup
	limiter := make(chan struct{}, 1000)
	numberOfRequests, start := 10_000_000, time.Now()
	for i := 1; i <= numberOfRequests; i++ {
		limiter <- struct{}{}
		wg.Add(1)
		go req(&wg, i, l, limiter)
	}

	wg.Wait()
	fmt.Printf("done in: %v\n", time.Since(start))
}
