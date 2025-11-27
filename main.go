package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/luis-octavius/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
}

func main() {
	godotenv.Load()

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)

	dbQueries := database.New(db)

	mux := http.NewServeMux()

	var apiCfg apiConfig
	apiCfg.queries = dbQueries

	// server config
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	appHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	// endpoints
	mux.Handle("/app/", appHandler)
	mux.Handle("GET /api/healthz", handlerHealthz())

	mux.Handle("GET /admin/metrics", apiCfg.handlerMetrics())
	mux.Handle("POST /admin/reset", apiCfg.handlerReset())
	mux.Handle("POST /api/validate_chirp", apiCfg.handlerValidateChirp())

	// ListenAndServe starts a server with an address and a handler
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("error listening on server: %w", err)
	}
}
