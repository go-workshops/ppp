package main

import (
	"log"
	"net/http"

	"github.com/go-workshops/ppp/cmd/simple-http/routes"
	"github.com/go-workshops/ppp/cmd/simple-http/services"
	"github.com/go-workshops/ppp/pkg/db"
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

	fileDB, err := db.OpenFile(".todos.db.json")
	if err != nil {
		log.Fatalln("could not open file database:", err)
	}

	todosSvc := services.NewTodo(fileDB)

	router := routes.NewRouter(routes.Config{
		TodosService: todosSvc,
	})
	log.Fatalln(http.ListenAndServe(":8080", router))
}
