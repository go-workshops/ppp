package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func line(level, msg string, fields map[string]interface{}) string {
	ts := time.Now().Format(time.RFC3339)
	logEntry := map[string]interface{}{
		"level": level,
		"msg":   msg,
		"ts":    ts,
	}
	for k, v := range fields {
		logEntry[k] = v
	}
	bs, _ := json.Marshal(logEntry)
	return string(bs)
}

func req(wg *sync.WaitGroup, id int, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }()
	log.Println(line("info", "request", map[string]any{"req_id": id}))
}

func main() {
	file, err := os.Create("std.log")
	if err != nil {
		log.Fatalf("could not create log file: %v", err)
	}
	defer func() { _ = file.Close() }()
	log.SetOutput(file)
	log.SetFlags(0)

	var wg sync.WaitGroup
	limiter := make(chan struct{}, 1000)
	numberOfRequests, start := 10_000_000, time.Now()
	for i := 1; i <= numberOfRequests; i++ {
		limiter <- struct{}{}
		wg.Add(1)
		go req(&wg, i, limiter)
	}
	wg.Wait()
	fmt.Printf("done in: %v\n", time.Since(start))
}
