package main

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/go-workshops/ppp/pkg/logging"
)

func handleRequest(wg *sync.WaitGroup, id int, logger *zap.Logger, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }() // Release the semaphore
	logger.Info("Handling request", zap.Int("request_id", id))
}

func main() {
	_ = logging.Init(logging.Config{
		LoggingLevel:  "debug",
		LoggingOutput: []string{"hello.log"},
	})
	logger := logging.GetLogger()

	var wg sync.WaitGroup
	numRequests := 10000000 // Simulate 1,000,000 incoming requests

	// Semaphore to limit concurrent goroutines
	limiter := make(chan struct{}, 1000) // Adjust the buffer size as needed

	start := time.Now()

	// Simulate handling a large number of concurrent requests
	for i := 1; i <= numRequests; i++ {
		limiter <- struct{}{} // Acquire a token
		wg.Add(1)
		go handleRequest(&wg, i, logger, limiter)
	}

	wg.Wait()

	fmt.Printf("All requests handled in %v", time.Since(start))
}
