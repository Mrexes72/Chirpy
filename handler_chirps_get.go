package main

import (
	"net/http"
	"sort"

	"github.com/Mrexes72/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp
	var err error

	authorIDStr := r.URL.Query().Get("author_id")

	if authorIDStr != "" {
		authorID, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, 400, "Invalid author ID", err)
			return
		}
		dbChirps, err = cfg.database.GetChirpsByAuthorID(r.Context(), authorID)
	} else {
		dbChirps, err = cfg.database.GetAllChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, 500, "Could not get chirps", err)
		return
	}

	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	sortOrder := r.URL.Query().Get("sort")

	if sortOrder == "" {
		sortOrder = "asc"
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, 400, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.database.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Could not get chirp", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, 200, chirp)
}
