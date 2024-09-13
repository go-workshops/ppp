package main

import (
	"log"
	"net/http"

	"github.com/go-workshops/ppp/cmd/simple-metrics/routes"
	"github.com/go-workshops/ppp/pkg/logging"
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes.NewRouter(),
	}
	log.Fatalln(srv.ListenAndServe())
}
