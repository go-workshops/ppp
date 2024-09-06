package main

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"buffered_zap.log"}
	cfg.Sampling = nil

	l, _ := cfg.Build()

	ws, _, err := zap.Open(cfg.OutputPaths...)
	bufferedWriteSyncer := &zapcore.BufferedWriteSyncer{
		WS:            ws,
		Size:          256 * 1024,
		FlushInterval: 2 * time.Second,
	}
	if err != nil {
		return nil, err
	}
	l = l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			bufferedWriteSyncer,
			cfg.Level,
		)
	}))

	return l, nil
}

func req(wg *sync.WaitGroup, id int, logger *zap.Logger, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }()
	logger.Info("request", zap.Int("req_id", id))
}

func main() {
	l, err := logger()
	if err != nil {
		panic(err)
	}
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
