package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/uLuKaiDev/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is DELETE
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	// Get the token from the request header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't check the token", err)
		return
	}

	// Validate the token and get the userID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Get the chirpID from the URL path
	idStr := r.PathValue("chirpID")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "Missing chirpID", nil)
		return
	}

	// Change the chirpID string to a UUID for database operations
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirpID", err)
		return
	}

	// Get the chirp struct from the database
	chirp, err := cfg.db.GetChirpsById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
		return
	}

	// Check if the chirp belongs to the user
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You are not allowed to delete this chirp", nil)
		return
	}

	// Delete the chirp from the database
	err = cfg.db.DeleteChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	// Respond with a success message
	respondWithJSON(w, http.StatusNoContent, nil)
}
