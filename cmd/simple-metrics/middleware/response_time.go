package middleware

import (
	"net/http"
	"time"

	"github.com/go-workshops/ppp/pkg/metrics"
)

var responseTimeHistogramMetric = metrics.HistogramVecWithBuckets(
	"http_response_time_ms",
	[]float64{50, 100, 200, 300, 400, 500, 1000},
	"Response time of the HTTP requests",
	"path",
)

func ResponseTime(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		duration := float64(time.Since(start).Milliseconds())
		responseTimeHistogramMetric.With(map[string]string{
			"path": r.URL.Path,
		}).Observe(duration)
	})
}
