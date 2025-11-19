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

	mux.Handle("/", http.FileServer(http.Dir(".")))

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("error listening on server: %w", err)
	}
}
