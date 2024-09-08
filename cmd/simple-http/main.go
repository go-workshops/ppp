package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-workshops/ppp/cmd/simple-http/routes"
	"github.com/go-workshops/ppp/cmd/simple-http/services"
	sharedContext "github.com/go-workshops/ppp/pkg/context"
	"github.com/go-workshops/ppp/pkg/db"
	"github.com/go-workshops/ppp/pkg/logging"
)

var (
	revision  string
	buildTime string
)

func main() {
	err := logging.Init(logging.Config{
		LoggingLevel:  "debug",
		LoggingOutput: []string{"stdout", "app.log"},
	})
	if err != nil {
		log.Fatalln("could not initialize logger:", err)
	}
	defer logging.Sync()

	fileDB, err := db.OpenFS(".db")
	if err != nil {
		log.Fatalln("could not open fs database:", err)
	}

	todosSvc := services.NewTodo(fileDB)

	router := routes.NewRouter(routes.Config{
		TodosService: todosSvc,
	})

	logger := logging.GetLogger()
	if revision != "" {
		logger = logger.With(zap.String("revision", revision))
	}
	if buildTime != "" {
		unixNano, _ := strconv.ParseInt(buildTime, 10, 64)
		if unixNano > 0 {
			logger = logger.With(zap.Time("build_time", time.Unix(0, unixNano)))
		}
	}
	ctx := sharedContext.WithLogger(context.Background(), logger)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		ErrorLog: logging.HTTPErrorLogger(),
	}
	log.Fatalln(srv.ListenAndServe())
}
