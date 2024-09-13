package routes

import (
	"net/http"

	"github.com/go-workshops/ppp/pkg/tracing"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/notify", instrument(notify(), "notify_user"))
	return mux
}

func instrument(h http.HandlerFunc, operation string) http.Handler {
	return tracing.InstrumentHTTP(tracing.HTTPMiddleware(h), operation)
}
