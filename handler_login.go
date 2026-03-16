package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mrexes72/Chirpy/internal/auth"
	"github.com/Mrexes72/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request", err)
		return
	}

	dbUser, err := cfg.database.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Invalid email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, 500, "Error checking password", err)
	}

	if !match {
		respondWithError(w, 401, "Invalid email or password", nil)
		return
	}

	const accessTokenExpiration = time.Hour
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.secretKey, accessTokenExpiration)
	if err != nil {
		respondWithError(w, 500, "Could not create access token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Could not create refresh token", err)
		return
	}

	const refreshTokenExpiration = 60 * 24 * time.Hour // 60 days
	_, err = cfg.database.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().Add(refreshTokenExpiration),
		RevokedAt: sql.NullTime{Valid: false},
	})
	if err != nil {
		respondWithError(w, 500, "Could not save refresh token", err)
		return
	}

	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	respondWithJSON(w, 200, response{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
		IsChirpyRed:  dbUser.IsRed,
	})
}
