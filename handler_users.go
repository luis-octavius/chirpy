package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/luis-octavius/chirpy/internal/auth"
	"github.com/luis-octavius/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsers() http.Handler {
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
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}

		writeJSON(w, http.StatusCreated, resp)
	})
}

func (cfg *apiConfig) handlerUserLogin() http.Handler {
	type reqParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Expires  int    `json:"expires_in_seconds"`
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
		var expires time.Duration
		if params.Expires >= 3600 || params.Expires == 0 {
			expires, err = time.ParseDuration("3600s")
		} else {
			expires, err = time.ParseDuration(strconv.Itoa(params.Expires) + "s")
		}

		token, err := auth.MakeJWT(user.ID, cfg.secret, expires)

		// create JSON answer
		resp := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		}

		writeJSON(w, http.StatusOK, resp)
	})
}
