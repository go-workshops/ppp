package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-workshops/ppp/playground/pprof-mem-leak-subtle/workers"
)

func main() {
	go func() {
		log.Fatalln(http.ListenAndServe(":8080", router()))
	}()

	go workers.Process(workers.Leak)
	go workers.Process(workers.NoLeak)
	time.Sleep(time.Hour)
}

func router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/pprof/", pprof.Index)
	mux.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/pprof/profile", pprof.Profile)
	mux.HandleFunc("/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/pprof/trace", pprof.Trace)
	mux.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/pprof/mutex", pprof.Handler("mutex"))
	mux.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/pprof/block", pprof.Handler("block"))
	return mux
}
