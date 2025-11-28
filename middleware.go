package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/luis-octavius/chirpy/internal/database"
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
			fmt.Println("is not dev")
			return
		}

		err := cfg.queries.DeleteAllUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("not deleted")
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func (cfg *apiConfig) handlerChirp() http.Handler {
	type CreateChirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type validateChirpResponse struct {
		Error       string `json:"error,omitempty"`
		Valid       bool   `json:"valid,omitempty"`
		CleanedBody string `json:"cleaned_body,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateChirpRequest

		// decoding POST body request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := validateChirpResponse{Error: "Something went wrong"}
			writeJSON(w, http.StatusInternalServerError, resp)
			return
		}

		chirpLength := 140

		if len(req.Body) > chirpLength {
			resp := validateChirpResponse{Error: "Chirp is too long"}
			writeJSON(w, http.StatusBadRequest, resp)
			return
		}

		filteredMessage := validateMessage(req.Body)

		chirp, err := cfg.queries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   filteredMessage,
			UserID: req.UserID,
		})
		if err != nil {
			log.Printf("error creating the chirp: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}

		writeJSON(w, http.StatusCreated, resp)
	})
}

func (cfg *apiConfig) handlerUsers() http.Handler {
	type validateParams struct {
		Email string `json:"email,omitempty"`
		Error string `json:"error,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req validateParams

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := validateParams{Error: "something went wrong with JSON decoding"}
			writeJSON(w, http.StatusInternalServerError, resp)
			return
		}

		user, err := cfg.queries.CreateUser(r.Context(), req.Email)
		if err != nil {
			resp := validateParams{Error: "something went wrong creating user in database"}
			writeJSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		writeJSON(w, http.StatusCreated, resp)

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
