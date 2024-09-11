package control_gc

import (
	"runtime/debug"
	"testing"
)

func BenchmarkMyFunction(b *testing.B) {
	debug.SetGCPercent(-1)
	defer debug.SetGCPercent(100)

	for i := 0; i < b.N; i++ {
		// MyFunction()
	}
}
