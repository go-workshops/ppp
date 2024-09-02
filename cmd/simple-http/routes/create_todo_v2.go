package routes

import (
	"encoding/json"
	"net/http"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

type createTodoRequestV2 struct {
	Title string `json:"title"`
}

func createTodoV2() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := sharedContext.Logger(r.Context())

		var req createTodoRequestV2
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			logger.Error("could not decode json request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Title == "" {
			logger.Error("could not create todo with missing title")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		logger.Info("successfully created todo")
		w.WriteHeader(http.StatusOK)
	}
}
