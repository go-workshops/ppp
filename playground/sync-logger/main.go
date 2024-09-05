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
	logEntry := map[string]interface{}{"ts": ts, "level": level, "msg": msg}
	for k, v := range fields {
		logEntry[k] = v
	}
	jsonLog, _ := json.Marshal(logEntry)
	return string(jsonLog)
}

func req(wg *sync.WaitGroup, id int, limiter chan struct{}) {
	defer wg.Done()
	defer func() { <-limiter }()
	log.Println(line("INFO", "request", map[string]any{"req_id": id}))
}

func main() {
	file, err := os.Create("std.log")
	if err != nil {
		log.Fatalf("Error creating log file: %v", err)
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
