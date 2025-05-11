package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/uLuKaiDev/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		log.Fatal("JWTSECRET must be set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dbQueries := database.New(dbConn)

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/login", cfg.handlerUsersLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerTokensRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerTokensRevoke)

	mux.HandleFunc("POST /api/chirps", cfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerChirpsGetByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerChirpsDelete)

	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", cfg.handlerUsersUpdate)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerPolkaWebhook)

	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	log.Printf("Serving on port %s\n", port)
	log.Printf("Open http://localhost:%s in your browser\n", port)

	// Listen and Serve has to be at the end of the main function
	// This ensures the server is running. The log.fatal is used only when an error occurs.
	log.Fatal(srv.ListenAndServe())
}
