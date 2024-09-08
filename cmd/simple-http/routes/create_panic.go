package routes

import (
	"net/http"
)

func createPanic() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("this will panic")
	}
}
