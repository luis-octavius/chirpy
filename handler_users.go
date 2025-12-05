package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/luis-octavius/chirpy/internal/auth"
	"github.com/luis-octavius/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser() http.Handler {
	type validateParams struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password"`
		Error    string `json:"error,omitempty"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req validateParams

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := validateParams{Error: "something went wrong with JSON decoding"}
			writeJSON(w, http.StatusBadRequest, resp)
			return
		}

		hashPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			resp := validateParams{Error: "error hashing the password"}
			writeJSON(w, http.StatusInternalServerError, resp)
			return
		}

		user, err := cfg.queries.CreateUser(r.Context(), database.CreateUserParams{
			Email:          req.Email,
			HashedPassword: hashPassword,
		})
		if err != nil {
			resp := validateParams{Error: "something went wrong creating user in database"}
			writeJSON(w, http.StatusInternalServerError, resp)
			return
		}

		resp := User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		}

		writeJSON(w, http.StatusCreated, resp)
	})
}

func (cfg *apiConfig) handlerUserLogin() http.Handler {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params reqParams

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// retrieves user from db - fails if user doesn't exist or use wrong password
		user, err := cfg.queries.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect email and password"))
			return
		}

		// check password hash against input password
		checkPassword, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
		if err != nil || !checkPassword {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Incorrect email or password"))
			return
		}

		// expires cannot be greather than 1 hour
		// default is 1 hour
		expiresAccToken := 1 * time.Hour
		refreshToken, _ := auth.MakeRefreshToken()

		token, err := auth.MakeJWT(user.ID, cfg.secret, expiresAccToken)

		newRefreshToken, err := cfg.queries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().AddDate(0, 0, 60),
		})

		// create JSON answer
		resp := User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: newRefreshToken.Token,
			IsChirpyRed:  user.IsChirpyRed,
		}

		writeJSON(w, http.StatusOK, resp)
	})
}

func (cfg *apiConfig) handlerUpdateUser() http.Handler {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params reqParams

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("error getting token from Authentication Header: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		validatedUser, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			log.Printf("error authenticating token: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := cfg.queries.GetUserByID(r.Context(), validatedUser)
		if err != nil {
			log.Printf("error getting user by ID: %v", err)
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("error hashing password: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = cfg.queries.UpdateUserEmailAndPass(r.Context(), database.UpdateUserEmailAndPassParams{
			HashedPassword: hashedPassword,
			Email:          params.Email,
			ID:             user.ID,
		})
		if err != nil {
			log.Printf("error updating user email and password: %v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		resp := User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       params.Email,
			IsChirpyRed: user.IsChirpyRed,
		}
		fmt.Println("resp: ", resp)

		writeJSON(w, http.StatusOK, resp)
	})
}

func (cfg *apiConfig) handlerUpgradeUser() http.Handler {

	type reqParams struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params reqParams

		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if params.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		userID, err := uuid.Parse(params.Data.UserID)
		err = cfg.queries.UpgradeUserByID(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return

	})
}
