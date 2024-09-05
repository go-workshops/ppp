package routes

import (
	"net/http"

	"github.com/go-workshops/ppp/cmd/simple-http/middleware"
)

type todosService interface {
	todoCreator
	todoUpdater
}

type Config struct {
	TodosService todosService
}

func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()

	// Don't mind my primitive routing, I know it's not REST like
	// I avoided using a router library for brevity and simplicity
	mux.HandleFunc("/v1/todos", createTodoV1())
	mux.HandleFunc("/v2/todos", createTodoV2())
	mux.HandleFunc("/v3/todos", createTodoV3(cfg.TodosService))
	mux.HandleFunc("/v4/todos", createTodoV4(cfg.TodosService))
	mux.HandleFunc("/v1/todos/update", updateTodoV1(cfg.TodosService))
	mux.HandleFunc("/v2/todos/update", updateTodoV2(cfg.TodosService))
	mux.HandleFunc("/v3/todos/update", updateTodoV3(cfg.TodosService))
	mux.HandleFunc("/v4/todos/update", updateTodoV4(cfg.TodosService))
	mux.HandleFunc("/v5/todos/update", updateTodoV5(cfg.TodosService))

	return middleware.New(
		mux,
		// middleware.RequestDumpV1,
		//middleware.RequestDumpV2,
		middleware.RequestDumpV3,
		// ...
	)
}
