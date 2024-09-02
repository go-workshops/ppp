package routes

import (
	"net/http"

	"github.com/go-workshops/ppp/cmd/simple-http/middleware"
)

type todosService interface {
	todoCreator
}

type Config struct {
	TodosService todosService
}

func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/todos", createTodoV1())
	mux.HandleFunc("/v2/todos", createTodoV2())
	mux.HandleFunc("/v3/todos", createTodoV3(cfg.TodosService))

	return middleware.New(
		mux,
		// middleware.RequestDumpV1,
		middleware.RequestDumpV2,
		// ...
	)
}
