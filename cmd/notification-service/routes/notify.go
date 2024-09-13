package routes

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
	"github.com/go-workshops/ppp/pkg/tracing"
)

func notify() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// omitting the service layer in here for brevity
		ctx := r.Context()
		userID := r.URL.Query().Get("user_id")
		logger := sharedContext.Logger(ctx).With(zap.String("user_id", userID))

		_, span := tracing.StartHTTP(ctx, "mailgun_service", "notify_user")
		defer span.End()
		span.SetAttributes(attribute.String("user_id", userID))

		// simulate sending the notification
		time.Sleep(time.Second)

		logger.Info("successfully notified user")
		w.WriteHeader(http.StatusOK)
	}
}
