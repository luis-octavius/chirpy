package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// middlewareMetricsInc is a wrapper that adds one to the count of
// the times that /app endpoint has been hit
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// handlerMetrics returns html code to render how many times
// /app endpoint has been hit
func (cfg *apiConfig) handlerMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		html := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	})
}

// handlerReset reset the count of times that app endpoint has been hit
// it also delete all users from database based on the platform that are
// developing the app
func (cfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Swap(0)

		if cfg.platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err := cfg.queries.DeleteAllUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

// writeJSON is a helper function that marshals data
// and write the same data into the ResponseWriter
//
// always returns JSON and write to the Header and
// the body of response the status of the request
func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

// validateMessage gets a string and replace
// all the ocurrences of prohibited words with four(*)
// then returns the filtered string
func validateMessage(s string) string {
	prohibitedWords := map[string]string{
		"kerfuffle": "kerfuffle",
		"sharbert":  "sherbert",
		"fornax":    "fornax",
	}

	splittedWords := strings.Split(s, " ")
	cleared := make([]string, len(splittedWords))

	for _, word := range splittedWords {
		loweredWord := strings.ToLower(word)
		_, ok := prohibitedWords[loweredWord]
		if ok {
			cleared = append(cleared, "****")
			continue
		}
		cleared = append(cleared, word)
	}

	return strings.TrimSpace(strings.Join(cleared, " "))
}
