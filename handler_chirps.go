package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/luis-octavius/chirpy/internal/auth"
	"github.com/luis-octavius/chirpy/internal/database"
)

// handlerAddChirps adds a chirp on the database
//
// Returns 400 if JSON decoding fails or chirp exceeds length limit
// Returns 500 if the chirp creation fails on the database
// Returns 201 with created chirp data on success
func (cfg *apiConfig) handlerAddChirps() http.Handler {
	type CreateChirpRequest struct {
		Body string `json:"body"`
	}

	type validateChirpResponse struct {
		Error       string `json:"error,omitempty"`
		Valid       bool   `json:"valid,omitempty"`
		CleanedBody string `json:"cleaned_body,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateChirpRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := validateChirpResponse{Error: "Something went wrong"}
			writeJSON(w, http.StatusBadRequest, resp)
			return
		}

		bearerToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("error getting the Bearer token: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(bearerToken, cfg.secret)
		if err != nil {
			log.Printf("error authorizing user: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
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
			UserID: userID,
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

// handlerGetAllChirps retrieves all chirps from the database
//
// Returns 500 if the chirps cannot be retrieved from database
// Returns 200 with all chirps data on success
func (cfg *apiConfig) handlerGetAllChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fetchedChirps, err := cfg.queries.GetAllChirps(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, fetchedChirps)
	})
}

// HandlerGetChirp returns a chirp based on a id path
//
// Returns 400 if the chirp ID cannot be parsed as UUID
// Returns 404 if no chirp exists with the given ID
// Returns 200 with chirp data on success
func (cfg *apiConfig) handlerGetChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirpID := r.PathValue("chirpID")
		parsedChirpID, err := uuid.Parse(chirpID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		chirp, err := cfg.queries.GetChirpByID(r.Context(), parsedChirpID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}

		writeJSON(w, http.StatusOK, chirp) // 200
	})
}

func (cfg *apiConfig) handlerDeleteChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirpID := r.PathValue("chirpID")
		parsedChirpID, err := uuid.Parse(chirpID)
		if err != nil {
			log.Printf("error parsing chirpID into UUID: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("error getting the token from header: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		validatedUser, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			log.Printf("error validating JWT token from user: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		chirp, err := cfg.queries.GetChirpByID(r.Context(), parsedChirpID)
		if err != nil {
			log.Printf("error getting chirp by ID: %v", err)
		}

		if chirp.UserID != validatedUser {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		err = cfg.queries.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
			ID:   parsedChirpID,
			ID_2: validatedUser,
		})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("chirp with id %v deleted successfully", chirpID)
		w.WriteHeader(http.StatusNoContent)
	})
}
