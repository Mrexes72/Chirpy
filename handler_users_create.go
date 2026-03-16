package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mrexes72/Chirpy/internal/auth"
	"github.com/Mrexes72/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Something went wrong", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "Could not hash password", err)
		return
	}

	dbUser, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		respondWithError(w, 500, "Could not create user", err)
		return
	}

	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsRed,
	}

	respondWithJSON(w, 201, user)
}
