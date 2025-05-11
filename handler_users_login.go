package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/uLuKaiDev/Chirpy/internal/auth"
	"github.com/uLuKaiDev/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	const accessTokenExpiry = time.Hour

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "missing password", nil)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get user", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", nil)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, accessTokenExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token", err)
		return
	}

	newToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: newToken.Token,
			IsChirpyRed:  user.IsChirpyRed.Bool,
		},
	})
}
