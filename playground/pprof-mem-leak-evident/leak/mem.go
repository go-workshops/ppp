package leak

import (
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

		go func(res *Resource) {
			for {
				time.Sleep(1 * time.Second)
				fmt.Println("Working on resource:", res.id)
				// Goroutine never exits, holding onto the `res` forever
			}
		}(r)

		time.Sleep(time.Second) // Simulate a new resource processing every second
	}
}
