package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
	"github.com/go-workshops/ppp/pkg/logging"
)

// go build -ldflags "-X 'main.revision=$(git rev-parse --short HEAD)' -X 'main.buildTime=$(date +%s000000000)'" -o bin/env-context playground/env-context/main.go
var (
	revision  string
	buildTime string
)

func main() {
	logger := logging.GetLogger().With(
		zap.String("revision", revision),
		zap.String("build_time", buildTime),
	)
	ctx := sharedContext.WithLogger(context.Background(), logger)

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("hello world")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	log.Fatalln(srv.ListenAndServe())
}
