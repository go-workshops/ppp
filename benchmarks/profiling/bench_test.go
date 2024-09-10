package profiling

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
)

func BenchmarkRouterProfilingDisabled(b *testing.B) {
	file, _ := os.Create("p1.log")
	defer func() { _ = file.Close() }()
	log.SetOutput(file)

	handler := router(false)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/test?req=%d", i), nil)
		handler.ServeHTTP(w, req)
	}
}

func BenchmarkRouterProfilingEnabled(b *testing.B) {
	file, _ := os.Create("p2.log")
	defer func() { _ = file.Close() }()
	log.SetOutput(file)

	handler := router(true)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/test?req=%d", i), nil)
		handler.ServeHTTP(w, req)
	}
}
