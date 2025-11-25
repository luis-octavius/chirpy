package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	mux := http.NewServeMux()

	var apiCfg apiConfig

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	appHandler := http.FileServer(http.Dir("."))

	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(appHandler)))
	mux.Handle("/healthz", handlerHealthz())
	mux.Handle("/metrics", apiCfg.handlerMetrics())
	mux.Handle("/reset", apiCfg.handlerReset())

	// ListenAndServe starts a server with an address and a handler
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("error listening on server: %w", err)
	}
}
