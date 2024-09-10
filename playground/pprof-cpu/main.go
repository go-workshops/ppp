package main

import (
	"log"
	"net/http"
	"net/http/pprof"
	"time"
)

func main() {
	go func() {
		log.Fatalln(http.ListenAndServe(":8080", pprofRouter()))
	}()

	for i := 0; i < 100; i++ {
		go work()
		time.Sleep(time.Second)
	}

	time.Sleep(10 * time.Minute)
}

func work() {
	for {
		getHot()
		time.Sleep(100 * time.Millisecond)
	}
}

func getHot() {
	sum := 0
	for i := 0; i < 1e6; i++ {
		sum += i
	}
}

func pprofRouter() http.Handler {
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
