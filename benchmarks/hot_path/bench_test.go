package main

import (
	"fmt"
	"sync"
	"testing"

	"go.uber.org/zap"
)

func BenchmarkTx1(b *testing.B) {
	var wg sync.WaitGroup
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"zap1.log"}
	logger, _ := cfg.Build()
	session := postgres()
	limiter := make(chan struct{}, 100)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		limiter <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-limiter }()
			tx1(logger, session, &wg, limiter, req{
				Author: fmt.Sprintf("Author %d", i),
				Book:   fmt.Sprintf("Book %d", i),
			})
		}(i)
	}
	wg.Wait()
}

func BenchmarkTx2(b *testing.B) {
	var wg sync.WaitGroup
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"zap2.log"}
	logger, _ := cfg.Build()
	session := postgres()
	limiter := make(chan struct{}, 100)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		limiter <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-limiter }()
			tx2(logger, session, &wg, limiter, req{
				Author: fmt.Sprintf("Author %d", i),
				Book:   fmt.Sprintf("Book %d", i),
			})
		}(i)
	}
	wg.Wait()
}
