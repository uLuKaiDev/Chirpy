package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/uLuKaiDev/Chirpy/internal/auth"
	"github.com/uLuKaiDev/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerTokensRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: token,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
