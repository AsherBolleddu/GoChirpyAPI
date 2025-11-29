package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AsherBolleddu/GoChirpyAPI/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API key", err)
		return
	}

	if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	params := parameters{}
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err = cfg.db.UpdateChirpyRed(r.Context(), params.Data.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusNotFound, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
