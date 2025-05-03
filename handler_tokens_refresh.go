package main

import (
	"net/http"
	"time"

	"github.com/uLuKaiDev/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerTokensRefresh(w http.ResponseWriter, r *http.Request) {
	// This functions checks the header for a token, which is just a string.
	// After that it checks if it's in the database or if it's revoked or expired.
	// If this checks out, a new JWT token is created and returned.
	// There is no new refresh token created, so the old one will be limited by it's validity.

	const duration = time.Hour

	// Checks the header for a bearer token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", err)
		return
	}

	// Checks the database for the refresh token specified in the header above
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Checks if the refresh token is revoked or expired. This checks the DATABASE token, which has
	// additional fields, whereas the normal "token" is just a string.
	if refreshToken.RevokedAt.Valid || refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired or revoked", nil)
	}

	// Looks up the user in the database using the refreshToken, it uses the .Token field of
	// the refreshToken struct, which is the string stored in the database.
	user, err := cfg.db.GetUserByRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// It issues a new JWT token and returns it.
	newToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{
		"token": newToken,
	})

}
