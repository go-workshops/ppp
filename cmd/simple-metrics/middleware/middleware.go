package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func New(h http.Handler, middlewares ...Middleware) http.Handler {
	// Apply middlewares in reverse order (last middleware first)
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
