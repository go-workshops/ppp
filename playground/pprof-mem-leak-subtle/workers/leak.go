package workers

import (
	"context"
	"fmt"
	"time"
)

func Leak(_ context.Context, res *Resource) {
	for {
		time.Sleep(time.Second)
		fmt.Println("Working on:", res.id)
	}
}
