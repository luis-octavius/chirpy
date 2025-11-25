package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(
			[]byte(fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())))
		if err != nil {
			fmt.Println("error writing to body", err)
		}
	})
}

func (cfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Swap(0)
	})
}
