package main

import (
	"net/http"
	"strconv"

	"github.com/closknight/Chirpy/internal/auth"
)

func (cfg apiConfig) HandleDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpID")
	id, err := strconv.Atoi(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id")
		return
	}
	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	jwtTokenString, err := auth.ParseBearer(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "no auth header found")
		return
	}
	idString, err := auth.ValidateJWT(jwtTokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not authenticate jwt")
		return
	}

	userId, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not convert user id from string to int")
		return
	}

	if chirp.AuthorID != userId {
		respondWithError(w, http.StatusForbidden, "user does not have permissions")
		return
	}

	err = cfg.DB.DeleteChirp(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
