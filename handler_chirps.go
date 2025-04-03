package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/willthefoollearn/chirpy/internal/auth"
	"github.com/willthefoollearn/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to retrieve token", err)
	}

	user, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token doesn't match user", err)
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to decode Chirp parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	splitBody := strings.Split(params.Body, " ")
	for i, word := range splitBody {
		low := strings.ToLower(word)
		if low == "kerfuffle" || low == "sharbert" || low == "fornax" {
			splitBody[i] = "****"
		}
	}
	cleanedBody := strings.Join(splitBody, " ")
	var args = database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: user,
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), args)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to create Chirp", err)
		return
	}

	respondWithJson(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, req *http.Request) {
	author_id := req.URL.Query().Get("author_id")
	sortMethod := req.URL.Query().Get("sort")
	if sortMethod != "desc" {
		sortMethod = "asc"
	}

	var chirps []database.Chirp
	var err error

	if author_id == "" {
		chirps, err = cfg.db.GetChirps(req.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Can't retrieve Chirps", nil)
		}
	} else {
		user, err := uuid.Parse(author_id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Can't retrieve user", nil)
		}

		chirps, err = cfg.db.GetChirpsByUser(req.Context(), user)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Can't retrieve Chirps", nil)
		}
	}

	convertedChirps := []Chirp{}

	for _, chirp := range chirps {
		convertedChirps = append(convertedChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.CreatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	if sortMethod == "desc" {
		sort.Slice(convertedChirps, func(i, j int) bool {
			comparison := convertedChirps[i].CreatedAt.Compare(convertedChirps[j].CreatedAt)

			return comparison == 1
		})
	}

	respondWithJson(w, http.StatusOK, convertedChirps)
}

func (cfg *apiConfig) handlerRetrieveChirp(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("chirpID")
	parsedUuid, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Can't parse path", nil)
	}

	chirp, err := cfg.db.RetrieveChirp(req.Context(), parsedUuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Can't retrieve Chirp", nil)
	}

	respondWithJson(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.CreatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	path := req.PathValue("chirpID")
	parsedUuid, err := uuid.Parse(path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Can't parse path", err)
	}

	chirp, err := cfg.db.RetrieveChirp(req.Context(), parsedUuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Can't retrieve Chirp", err)
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "BearerToken unable to be retrieved", err)
		return
	}

	user, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User can't be found with accessToken", err)
		return
	}

	if chirp.UserID == user {
		err = cfg.db.DeleteChirp(req.Context(), chirp.ID)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "No chirp to delete, ID didn't match", err)
			return
		}
	} else {
		respondWithError(w, http.StatusForbidden, "Incorrect user for Chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
