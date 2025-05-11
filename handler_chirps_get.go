package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "" {
		sortOrder = "asc"
	}

	// Check if the request had an optional user ID named author_id.
	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		// Validate the author ID
		userID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
		// Get all chirps by author ID
		chirps, err := cfg.db.GetChirpsByUserId(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
		response := make([]ChirpResponse, len(chirps))
		for i, chirp := range chirps {
			response[i] = ChirpResponse{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
		}
		if sortOrder == "desc" {
			// Reverse the order of the response slice
			sort.Slice(response, func(i, j int) bool {
				return response[i].CreatedAt.After(response[j].CreatedAt)
			})
		}

		respondWithJSON(w, http.StatusOK, response)
	}

	chirps, err := cfg.db.GetChirpsAsc(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	response := make([]ChirpResponse, len(chirps))
	for i, chirp := range chirps {
		response[i] = ChirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}

	if sortOrder == "desc" {
		// Reverse the order of the response slice
		sort.Slice(response, func(i, j int) bool {
			return response[i].CreatedAt.After(response[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, response)
}
