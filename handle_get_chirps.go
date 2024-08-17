package main

import (
	"net/http"
	"slices"
	"strconv"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

func (cfg *apiConfig) HandleGetChirps(w http.ResponseWriter, r *http.Request) {
	authorIDString := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	authorID := -1
	if len(authorIDString) > 0 {
		var err error
		authorID, err = strconv.Atoi((authorIDString))
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author id")
			return
		}
	}

	chirps, err := cfg.DB.GetChirps()

	if sortOrder == "desc" {
		slices.Reverse(chirps)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve chirps")
		return
	}
	results := []Chirp{}
	for _, chirp := range chirps {
		if authorID == -1 || authorID == chirp.AuthorID {
			results = append(results, Chirp{ID: chirp.Id, AuthorID: chirp.Id, Body: chirp.Body})
		}
	}
	respondWithJSON(w, http.StatusOK, results)
}

// func (db *DB) GetChirpsByAuthor(AuthorID int) ([]Chirp, error) {
// 	chirps, err := db.GetChirps()
// 	if err != nil {
// 		return []Chirp{}, err
// 	}
// 	results := []Chirp{}
// 	for _, chirp := range chirps {
// 		if chirp.AuthorID == AuthorID {
// 			results = append(results, chirp)
// 		}
// 	}
// 	return results, nil
// }

func (cfg *apiConfig) HandleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK,
			Chirp{
				ID:       chirp.Id,
				AuthorID: chirp.AuthorID,
				Body:     chirp.Body,
			})
	}
}
