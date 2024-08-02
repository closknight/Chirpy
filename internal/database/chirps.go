package database

import (
	"errors"
	"slices"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp, ok := dbs.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("resource does not exist")
	}
	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	dbs, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := make([]Chirp, 0, len(dbs.Chirps))
	for _, v := range dbs.Chirps {
		chirps = append(chirps, v)
	}
	slices.SortFunc(chirps, func(a, b Chirp) int { return a.Id - b.Id })
	return chirps, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbs.Chirps) + 1
	chirp := Chirp{
		Body: body,
		Id:   id,
	}

	dbs.Chirps[id] = chirp
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}
