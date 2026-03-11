package main

import (
	"net/http"
	"time"

	"github.com/Mrexes72/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized,
			"Missing or invalid authorization", err)
		return
	}

	user, err := cfg.database.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized,
			"Invalid or expired refresh token", err)
		return
	}

	const accessTokenExpiration = time.Hour
	accessToken, err := auth.MakeJWT(user.ID, cfg.secretKey, accessTokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Could not create access token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, 200, response{
		Token: accessToken,
	})
}
