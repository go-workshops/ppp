package main

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkUnixTime(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.EpochNanosTimeEncoder
	cfg.OutputPaths = []string{"zap.log"}
	logger, _ := cfg.Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("some message")
	}
}

func BenchmarkISO8601Time(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.OutputPaths = []string{"zap.log"}
	logger, _ := cfg.Build()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("some message")
	}
}
