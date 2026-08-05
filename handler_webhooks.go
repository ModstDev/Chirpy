package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ModstDev/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API key doesn't match", err)
		return
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
