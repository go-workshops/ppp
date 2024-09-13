package routes

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

type userRegisterer interface {
	Register(context.Context) (string, error)
}

type userNotifier interface {
	Notify(ctx context.Context, userID string) error
}

func register(registerer userRegisterer, notifier userNotifier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := sharedContext.Logger(ctx)

		userID, err := registerer.Register(ctx)
		if err != nil {
			logger.Error("could not register user", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = notifier.Notify(ctx, userID); err != nil {
			logger.Error("could not notify user", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("user successfully registered")
		w.WriteHeader(http.StatusOK)
	}
}
