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
	for i := 0; i < 1000; i++ {
		r := &Resource{
			id:   i,
			data: make([]byte, 10*1024*1024), // 10MiB
		}
		go func(res *Resource) {
			for {
				time.Sleep(1 * time.Second)
				fmt.Println("Working on:", res.id)
			}
		}(r)

		time.Sleep(time.Second)
	}
}
