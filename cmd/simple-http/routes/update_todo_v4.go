package routes

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-workshops/ppp/cmd/simple-http/models"
	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

func updateTodoV4(svc todoUpdater) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := sharedContext.Logger(ctx)

		var req updateTodoRequestV4
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("could not decode json request body", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !req.Validate(r, w) {
			return
		}

		todo := models.Todo{ID: req.ID, Title: req.Title, Description: req.Description}
		if err := svc.UpdateTodo(ctx, todo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("successfully updated todo", zap.String("id", todo.ID))
		w.WriteHeader(http.StatusOK)
	}
}

type updateTodoRequestV4 struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (req updateTodoRequestV4) Validate(r *http.Request, w http.ResponseWriter) bool {
	logger := sharedContext.Logger(r.Context())
	if req.ID == "" {
		logger.Warn("could not update todo with missing id")
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if req.Title == "" && req.Description == "" {
		logger.Warn("could not update todo with missing title and description")
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	return true
}
