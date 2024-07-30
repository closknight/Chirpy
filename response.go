package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
