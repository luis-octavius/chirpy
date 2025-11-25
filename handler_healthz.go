package main

import (
	"net/http"
)

func handlerHealthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(200)

		_, _ = w.Write([]byte("OK"))
	})
}
