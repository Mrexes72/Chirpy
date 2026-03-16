package main

import (
	"encoding/json"
	"net/http"

	"github.com/Mrexes72/Chirpy/internal/auth"
	"github.com/Mrexes72/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Invalid token", err)
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Failed to hash password", err)
		return
	}

	updatedUser, err := cfg.database.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, 500, "Failed to update user", err)
		return
	}

	user := User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsRed,
	}

	respondWithJSON(w, 200, user)
}
