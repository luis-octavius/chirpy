package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

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

func (cfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Swap(0)
	})
}

func (cfg *apiConfig) handlerValidateChirp() http.Handler {
	type validateChirpParams struct {
		Body string `json:"body"`
	}

	type validateChirpResponse struct {
		Error       string `json:"error,omitempty"`
		Valid       bool   `json:"valid,omitempty"`
		CleanedBody string `json:"cleaned_body,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params validateChirpParams

		// decoding POST body request
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			resp := validateChirpResponse{Error: "Something went wrong"}
			writeJSON(w, http.StatusInternalServerError, resp)
			return
		}

		if len(params.Body) > 140 {
			resp := validateChirpResponse{Error: "Chirp is too long"}
			writeJSON(w, http.StatusBadRequest, resp)
			return
		}

		resp := validateChirpResponse{CleanedBody: validateMessage(params.Body)}
		writeJSON(w, http.StatusOK, resp)
	})
}

// writeJSON is a helper function that marshal data
// and write the data into the ResponseWriter
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
