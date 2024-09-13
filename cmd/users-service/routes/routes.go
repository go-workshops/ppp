package routes

import (
	"net/http"

	"github.com/go-workshops/ppp/pkg/tracing"
)

type UsersService interface {
	userRegisterer
}

type NotificationClient interface {
	userNotifier
}

type Config struct {
	UsersService
	NotificationClient
}

func NewRouter(cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/register", instrument(register(cfg.UsersService, cfg.NotificationClient), "register_user"))
	return mux
}

func instrument(h http.HandlerFunc, operation string) http.Handler {
	return tracing.InstrumentHTTP(tracing.HTTPMiddleware(h), operation)
}
