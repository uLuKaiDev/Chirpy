package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uLuKaiDev/Chirpy/internal/auth"
	"github.com/uLuKaiDev/Chirpy/internal/database"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't check the token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	} else if len(params.Body) == 0 {
		respondWithError(w, http.StatusBadRequest, "Chirp is empty", err)
		return
	}

	params.Body = badWordReplacement(params.Body)
	arg := database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), arg)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	response := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, response)
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
