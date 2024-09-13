package routes

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/go-workshops/ppp/cmd/simple-metrics/middleware"
	"github.com/go-workshops/ppp/pkg/metrics"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	rand.New(rand.NewSource(time.Now().UnixNano()))

	mux.Handle("/metrics", metrics.PrometheusHandler())
	mux.HandleFunc("/m1", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Int63n(900)) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/m2", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Int63n(900)) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})

	return middleware.New(
		mux,
		middleware.ResponseTime,
	)
}
