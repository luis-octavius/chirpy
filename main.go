package main 

import (
	"net/http"
	"log"
)

func main () {

	mux := http.NewServeMux() 

	server := http.Server{
		Addr: ":8080", 
		Handler: mux,
	}
	
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.Handle("/healthz/", http.StripPrefix("/healthz/", handlerHealthz()))


	// ListenAndServe starts a server with an address and a handler
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("error listening on server: %w", err)
	}
}

func handlerHealthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")	

		w.WriteHeader(200)

		_, _ = w.Write([]byte("OK"))
	})
}
