package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

// Recovery recovers the application from a potential/accidental panic and returns generic 500 response error.
func Recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := sharedContext.Logger(ctx)
		defer func() {
			err := panicToError(recover())
			if err != nil {
				logger.Error("unexpected panic", zap.Error(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func panicToError(value any) error {
	var panicErr error
	switch e := value.(type) {
	case nil:
		return nil
	case string:
		if e == "" {
			return nil
		}
		panicErr = errors.New(e)
	case error:
		panicErr = e
	default:
		panicErr = fmt.Errorf("unknown panic: %v", e)
	}
	return panicErr
}
