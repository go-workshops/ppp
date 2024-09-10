package workers

import (
	"context"
	"time"
)

type Resource struct {
	id   int
	data []byte
}

func Process(worker func(context.Context, *Resource)) {
	for i := 0; i < 1000; i++ {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		r := &Resource{
			id:   i,
			data: make([]byte, 10*1024*1024), // 10MiB
		}

		go func(ctx context.Context, r *Resource) {
			defer cancel()
			worker(ctx, r)
		}(ctx, r)
		time.Sleep(time.Second)
	}
}
