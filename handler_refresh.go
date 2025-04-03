package main

import (
	"net/http"
	"time"

	"github.com/willthefoollearn/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "BearerToken unable to be retrieved", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "RefreshToken can't be found", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInsufficientStorage, "Couldn't make access token", err)
	}

	respondWithJson(w, http.StatusOK, Response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "BearerToken unable to be retrieved", err)
		return
	}

	_, err = cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "RefreshToken can't be found to revoke", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "RefreshToken can't be found to revoke", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
