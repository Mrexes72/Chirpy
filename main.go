package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Mrexes72/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
	platform       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DB_URL: %s", dbURL)
	log.Printf("Platform: %s", platform)
	dbQueries := database.New(db)

	log.Printf("Connecting to database at %s", dbURL)

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		database:       dbQueries,
		platform:       platform,
	}

	serverMux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))
	serverMux.Handle("/app/", fsHandler)

	serverMux.HandleFunc("GET /admin/healthz", handlerReadiness)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

	serverMux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.resetUsersHandler)

	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.getAllChirpsHandler)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpByIDHandler)

	serverMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: serverMux,
	}

	log.Printf("Serving on port %s\n", port)
	server.ListenAndServe()
}
