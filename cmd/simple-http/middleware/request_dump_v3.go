package middleware

import (
	"net/http"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

func RequestDumpV3(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := sharedContext.Logger(ctx).With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		h.ServeHTTP(w, r.WithContext(sharedContext.WithLogger(ctx, logger)))
	})
}
