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
	for i := 0; i < 1000; i++ {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		r := &Resource{
			id:   i,
			data: make([]byte, 10*1024*1024), // 10MiB
		}
		go func(res *Resource, ctx context.Context) {
			defer cancel()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("Done with:", res.id)
					return
				default:
					time.Sleep(1 * time.Second)
					fmt.Println("Working on:", res.id)
				}
			}
		}(r, ctx)
		time.Sleep(time.Second)
	}
}
