package main

import (
	"log"
	"net/http"
	"time"

	"github.com/luis-octavius/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken() http.Handler {
	type respToken struct {
		Token string `json:"token"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := cfg.queries.GetUserByRefreshToken(r.Context(), refreshToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		accToken, err := auth.MakeJWT(user.ID, cfg.secret, 1*time.Hour)
		if err != nil {
			log.Printf("error creating JWT access token: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := respToken{
			Token: accToken,
		}

		writeJSON(w, http.StatusOK, resp)
	})
}

func (cfg *apiConfig) handlerRevokeToken() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)

		err = cfg.queries.RevokeRefreshToken(r.Context(), token)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

}
