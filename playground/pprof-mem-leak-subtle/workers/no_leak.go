package workers

import (
	"context"
	"fmt"
	"time"
)

func NoLeak(ctx context.Context, res *Resource) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Done processing resource:", res.id)
			return // Exit the goroutine when context is done
		default:
			time.Sleep(1 * time.Second)
			fmt.Println("Working on resource:", res.id)
		}
	}
}
