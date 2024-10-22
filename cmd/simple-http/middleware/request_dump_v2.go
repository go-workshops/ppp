package middleware

import (
	"net/http"

	"go.uber.org/zap"

	sharedContext "github.com/go-workshops/ppp/pkg/context"
)

func RequestDumpV2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := sharedContext.Logger(r.Context())

		logger.Debug(
			"request dump",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)

		h.ServeHTTP(w, r)
	})
}
