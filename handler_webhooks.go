package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/willthefoollearn/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			User_ID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "apiKey not found", err)
		return
	}

	if apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Invalid apiKey", nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode body", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.RedUpgrade(req.Context(), params.Data.User_ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to find user", err)
	}

	w.WriteHeader(http.StatusNoContent)

}
