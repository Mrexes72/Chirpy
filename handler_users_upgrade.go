package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Mrexes72/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	polkaKey := os.Getenv("POLKA_KEY")
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid header", err)
		return
	}

	if apiKey != polkaKey {
		respondWithError(w, 401, "Unauthorized", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Invalid request body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.database.UpgradeUserToRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "User not found", err)
		return
	}

	w.WriteHeader(204)
}
