package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mrexes72/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
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

	const defaultExpiration = time.Hour
	const maxExpiration = time.Hour

	expiresIn := defaultExpiration
	if params.ExpiresInSeconds != nil {
		requestedExpiration := time.Duration(*params.ExpiresInSeconds) * time.Second

		if requestedExpiration > maxExpiration {
			expiresIn = maxExpiration
		} else {
			expiresIn = requestedExpiration
		}
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secretKey, expiresIn)
	if err != nil {
		respondWithError(w, 500, "Could not create token", err)
		return
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	respondWithJSON(w, 200, response{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	})
}
