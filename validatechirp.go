package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type resp struct {
		CleanedBody string `json:"cleaned_body"`
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
	respondWithJSON(w, http.StatusOK, resp{CleanedBody: RemoveProfanity(chirp.Body)})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResp struct {
		Error string `json:"error"`
	}
	if code >= 500 {
		log.Printf("responding with error %d: %s", code, msg)
	}

	respondWithJSON(w, code, errorResp{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshaling: %s\n", dat)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func Contains[T comparable](lst []T, item T) bool {
	for _, i := range lst {
		if i == item {
			return true
		}
	}
	return false
}

func RemoveProfanity(msg string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	clean_msg := []string{}
	for _, word := range strings.Split(msg, " ") {
		if Contains(profanity, strings.ToLower(word)) {
			clean_msg = append(clean_msg, "****")
		} else {
			clean_msg = append(clean_msg, word)
		}
	}
	return strings.Join(clean_msg, " ")
}
