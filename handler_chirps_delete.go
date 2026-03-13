package main

import (
	"net/http"

	"github.com/Mrexes72/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid authorization header", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, 401, "Invalid token", err)
		return
	}

	chirpIDStr := r.PathValue("chirpID")
	if chirpIDStr == "" {
		respondWithError(w, 400, "Missing chirp ID", nil)
		return
	}

	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, 400, "Invalid Chirp ID format", err)
	}

	chirp, err := cfg.database.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Chirp not found", err)
	}

	if chirp.UserID != userID {
		respondWithError(w, 403, "You don't have permission to delete this chirp", err)
		return
	}

	err = cfg.database.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 500, "Could not delete chirp", err)
		return
	}

	w.WriteHeader(204)

}
