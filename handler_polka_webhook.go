package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/uLuKaiDev/Chirpy/internal/auth"
	"github.com/uLuKaiDev/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Check if the request is from Polka by checking if there is an API key
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error retrieving API key", err)
		return
	}

	// Check the API key to the stored key
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", nil)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		}
	}
	// Decode request body into an empty parameters struct
	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	// Validate the event before proceeding
	if params.Event == "" || params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Validate the user_id. This changes the type from string to uuid.UUID
	// and checks if the user is in the database.

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id", err)
		return
	}

	user, err := cfg.db.GetUserById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user", err)
		return
	}

	// Update the user in the database
	_, err = cfg.db.UpdateUserRed(r.Context(), database.UpdateUserRedParams{
		ID:          user.ID,
		IsChirpyRed: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// Respond with a success message which is an empty response body
	w.WriteHeader(http.StatusNoContent)
}
