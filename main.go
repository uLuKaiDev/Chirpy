package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()
	cfg := &apiConfig{}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	//fileserver endpoints
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	//api endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("/api/validate_chirp", handlerChirpValidate)

	//admin endpoints
	mux.HandleFunc("GET /admin/metrics", cfg.handleMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetrics)

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Printf("Open http://localhost:%s in your browser\n", port)

	// Listen and Serve has to be at the end of the main function
	// This ensures the server is running. The log.fatal is used only when an error occurs.
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	output := fmt.Sprintf(`
<html>
  	<body>
    	<h1>Welcome, Chirpy Admin</h1>
    	<p>Chirpy has been visited %d times!</p>
  	</body>
</html>
`, cfg.fileserverHits.Load())

	w.Write([]byte(output))
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type chirpRequest struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Valid bool   `json:"valid"`
		Error string `json:"error,omitempty"`
	}

	var chirp chirpRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding request: %s", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var resp returnVals
	if len(chirp.Body) > 140 {
		resp.Valid = false
		resp.Error = "Chirp is too long"
		w.Header().Set("Content-Type", "application/json: charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		resp.Valid = true
		resp.Error = ""
		w.Header().Set("Content-Type", "application/json: charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}

	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshaling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(dat)
}
