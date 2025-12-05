package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/luis-octavius/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
	platform       string
	secret         string
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {
	godotenv.Load()

	// get the url of database from .env
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	// open the connection with the database url
	db, err := sql.Open("postgres", dbUrl)

	// initialize the holding of all queries made with sqlc
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	var apiCfg apiConfig
	apiCfg.queries = dbQueries
	apiCfg.platform = platform
	apiCfg.secret = os.Getenv("SECRET")

	// server config
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	appHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// endpoints
	mux.Handle("/app/", appHandler)
	mux.Handle("GET /api/healthz", handlerHealthz())

	// admin endpoints
	mux.Handle("GET /admin/metrics", apiCfg.handlerMetrics())
	mux.Handle("POST /admin/reset", apiCfg.handlerReset())

	// users endpoints
	mux.Handle("POST /api/users", apiCfg.handlerCreateUser())
	mux.Handle("POST /api/login", apiCfg.handlerUserLogin())
	mux.Handle("PUT /api/users", apiCfg.handlerUpdateUser())
	mux.Handle("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser())

	// chirps endpoints
	mux.Handle("GET /api/chirps", apiCfg.handlerGetAllChirps())
	mux.Handle("POST /api/chirps", apiCfg.handlerAddChirps())
	mux.Handle("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp())
	mux.Handle("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp())

	// token endpoints
	mux.Handle("POST /api/refresh", apiCfg.handlerRefreshToken())
	mux.Handle("POST /api/revoke", apiCfg.handlerRevokeToken())

	// ListenAndServe starts a server with an address and a handler
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("error listening on server: %w", err)
	}
}
