package profiling

import (
	"log"
	"net/http"
	"net/http/pprof"
)

func router(withProfiling bool) http.Handler {
	mux := http.NewServeMux()

	// publicly exposed routes available on the ingress controller
	mux.HandleFunc("/v1/test", func(w http.ResponseWriter, r *http.Request) {
		req := r.URL.Query().Get("req")
		log.Printf("test %s\n", req)
	})

	if !withProfiling {
		return mux
	}

	// internal routes only available via port forwarding
	mux.HandleFunc("pprof/", pprof.Index)
	mux.HandleFunc("pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("pprof/profile", pprof.Profile)
	mux.HandleFunc("pprof/symbol", pprof.Symbol)
	mux.HandleFunc("pprof/trace", pprof.Trace)
	mux.Handle("pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("pprof/heap", pprof.Handler("heap"))
	mux.Handle("pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("pprof/block", pprof.Handler("block"))

	return mux
}
