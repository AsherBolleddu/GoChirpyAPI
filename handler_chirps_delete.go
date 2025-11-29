package main

import (
	"net/http"

	"github.com/AsherBolleddu/GoChirpyAPI/internal/auth"
	"github.com/AsherBolleddu/GoChirpyAPI/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can't delete this chirp", err)
		return
	}

	if err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirp.ID,
		UserID: userID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
