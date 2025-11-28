package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/luis-octavius/chirpy/internal/database"
)

func (cfg *apiConfig) handlerAddChirps() http.Handler {
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

		// the size of the body cannot be greater than the size of a chirp
		if len(req.Body) > chirpLength {
			resp := validateChirpResponse{Error: "Chirp is too long"}
			writeJSON(w, http.StatusBadRequest, resp)
			return
		}

		// filter message to block prohibited words
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

func (cfg *apiConfig) handlerGetAllChirps() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetchedChirps, err := cfg.queries.GetAllChirps(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, fetchedChirps)
	})
}

func (cfg *apiConfig) handlerGetChirp() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("chirpID")
		parsedID, err := uuid.Parse(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}

		chirp, err := cfg.queries.GetChirpByID(r.Context(), parsedID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}

		writeJSON(w, http.StatusOK, chirp) // 200
	})
}
