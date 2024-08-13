package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/closknight/Chirpy/internal/auth"
)

const MAX_CHIRP_LENGTH = 140

func (cfg *apiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	jwtToken, err := auth.ParseBearer(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not retrieve auth token")
		return
	}

	idString, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized user")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse token id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err = decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could't decode request")
		return
	}

	if len(chirp.Body) > MAX_CHIRP_LENGTH {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	newChirp, err := cfg.DB.CreateChirp(id, sanitize(chirp.Body))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
	} else {
		respondWithJSON(w, http.StatusCreated, newChirp)
	}
}

func sanitize(msg string) string {
	profanity := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	clean_msg := []string{}
	for _, word := range strings.Split(msg, " ") {
		if _, ok := profanity[word]; ok {
			clean_msg = append(clean_msg, "****")
		} else {
			clean_msg = append(clean_msg, word)
		}
	}
	return strings.Join(clean_msg, " ")
}
