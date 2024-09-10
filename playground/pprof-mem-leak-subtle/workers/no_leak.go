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
			fmt.Println("Done with:", res.id)
			return
		default:
			time.Sleep(1 * time.Second)
			fmt.Println("Working on:", res.id)
		}
	}
}
