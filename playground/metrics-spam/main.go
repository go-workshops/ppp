package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	limiter := make(chan struct{}, 30)
	numberOfRequests := 10000
	for i := 1; i <= numberOfRequests; i++ {
		limiter <- struct{}{}
		wg.Add(2)
		go reqM1(&wg, limiter)
		go reqM2(&wg, limiter)
		time.Sleep(500 * time.Millisecond)
	}

	wg.Wait()
}

func reqM1(wg *sync.WaitGroup, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }()
	res, err := http.Get("http://localhost:8080/m1")
	if err != nil {
		log.Fatalf("could not send request: %v", err)
	}
	defer func() { _ = res.Body.Close() }()
}

func reqM2(wg *sync.WaitGroup, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }()
	res, err := http.Get("http://localhost:8080/m2")
	if err != nil {
		log.Fatalf("could not send request: %v", err)
	}
	defer func() { _ = res.Body.Close() }()
}
