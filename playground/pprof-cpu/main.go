package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

func main() {
	// Start the pprof HTTP server
	go func() {
		log.Println(http.ListenAndServe(":8080", router()))
	}()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			performWork(id)
		}(i)

		time.Sleep(10 * time.Millisecond) // Simulate some delay between goroutine starts
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Simulate running indefinitely
	time.Sleep(10 * time.Minute)
}

func performWork(id int) {
	for {
		// Simulate CPU-bound work
		busyWork()
		time.Sleep(100 * time.Millisecond) // Sleep to simulate periodic work
	}
}

func busyWork() {
	sum := 0
	for i := 0; i < 1e6; i++ { // Intensive CPU computation
		sum += i
	}
}

func router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/pprof/", pprof.Index)
	mux.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/pprof/profile", pprof.Profile)
	mux.HandleFunc("/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/pprof/trace", pprof.Trace)
	mux.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/pprof/block", pprof.Handler("block"))
	return mux
}
