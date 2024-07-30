package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const MAX_CHIRP_LENGTH = 140

func (cfg *apiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := parameters{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could't decode request")
		return
	}

	if len(chirp.Body) > MAX_CHIRP_LENGTH {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	newChirp, err := cfg.DB.CreateChirp(sanitize(chirp.Body))
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
