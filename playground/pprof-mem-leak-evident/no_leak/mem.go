package no_leak

import (
	"context"
	"fmt"
	"time"
)

type Resource struct {
	id   int
	data []byte
}

func Process() {
	// Simulate a streaming service
	for i := 0; i < 1000; i++ {
		r := &Resource{
			id:   i,
			data: make([]byte, 10*1024*1024), // 10MB
		}

		// Start a background goroutine that holds onto the resource, but will terminate after timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		go func(res *Resource, ctx context.Context) {
			defer cancel()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Done processing resource:", res.id)
					return // Exit the goroutine when context is done
				default:
					// Do some background work
					time.Sleep(1 * time.Second)
					fmt.Println("Working on resource:", res.id)
				}
			}
		}(r, ctx)

		time.Sleep(time.Second) // Simulate a new resource processing every second
	}
}
