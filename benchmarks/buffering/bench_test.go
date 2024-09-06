package main

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkWithBuffering(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.EpochNanosTimeEncoder
	cfg.OutputPaths = []string{"zap.log"}
	cfg.Sampling = nil
	logger, _ := cfg.Build()
	ws, _, _ := zap.Open(cfg.OutputPaths...)
	bufferedWriteSyncer := &zapcore.BufferedWriteSyncer{
		WS:            ws,
		Size:          256 * 1024,
		FlushInterval: 2 * time.Second,
	}
	logger = logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			bufferedWriteSyncer,
			cfg.Level,
		)
	}))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("some message")
	}
}

func BenchmarkWithoutBuffering(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"zap.log"}
	cfg.Sampling = nil
	logger, _ := cfg.Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("some message")
	}
}
