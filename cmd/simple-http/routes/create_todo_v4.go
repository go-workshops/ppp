package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-workshops/ppp/cmd/simple-http/models"
	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

type createTodoRequestV4 struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createTodoV4(svc todoCreator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := sharedContext.Logger(ctx)

		var req createTodoRequestV4
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("could not decode json request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Title == "" {
			logger.Warn("could not create todo with missing title")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Description == "" {
			logger.Warn("could not create todo with missing description")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todo := models.Todo{
			ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
			Title:       req.Title,
			Description: req.Description,
		}
		if err := svc.CreateTodo(ctx, todo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("successfully created todo")
		w.WriteHeader(http.StatusOK)
	}
}
