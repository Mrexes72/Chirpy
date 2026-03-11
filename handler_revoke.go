package main

import (
	"net/http"

	"github.com/Mrexes72/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized,
			"Missing or invalid authorization", err)
		return
	}

	err = cfg.database.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Could not revoke refresh token", err)
		return
	}

	w.WriteHeader(204)
}
