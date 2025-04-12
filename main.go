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
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dbQueries := database.New(dbConn)

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	//fileserver endpoints
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	//api endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpValidate)

	//admin endpoints
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetrics)

	log.Printf("Serving on port %s\n", port)
	log.Printf("Open http://localhost:%s in your browser\n", port)

	// Listen and Serve has to be at the end of the main function
	// This ensures the server is running. The log.fatal is used only when an error occurs.
	log.Fatal(srv.ListenAndServe())
}
