package workers

import (
	"context"
	"fmt"
	"time"
)

func Leak(_ context.Context, res *Resource) {
	for {
		time.Sleep(1 * time.Second)
		fmt.Println("Working on resource:", res.id)
		// Goroutine never exits, holding onto the `res` forever
	}
}
