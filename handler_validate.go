package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type chirpRequest struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Cleaned_body string `json:"cleaned_body,omitempty"`
	}

	var chirp chirpRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	var resp returnVals
	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	} else if len(chirp.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp is empty", err)
		return
	}
	chirp.Body = badWordReplacement(chirp.Body)
	resp.Cleaned_body = chirp.Body
	respondWithJSON(w, http.StatusOK, resp)

	// Uncomment the following line to simulate an internal server error
	// respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	// This line should never be reached
}

func badWordReplacement(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	wordsInBody := strings.Split(body, " ")
	for i, word := range wordsInBody {
		for _, badWord := range badWords {
			if strings.EqualFold(strings.ToLower(word), badWord) {
				wordsInBody[i] = "****"
			}
		}
	}
	return strings.Join(wordsInBody, " ")
}
