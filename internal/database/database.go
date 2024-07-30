package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"slices"
	"sync"
)

var chirps_count = 1

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	database := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := database.ensureDB()
	if err != nil {
		return &DB{}, err
	}
	return &database, nil
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
	chirp := Chirp{
		Body: body,
		Id:   chirps_count,
	}

	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	dbs.Chirps[len(dbs.Chirps)+1] = chirp
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}
	chirps_count++
	return chirp, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := os.ReadFile(db.path)

	if errors.Is(err, fs.ErrNotExist) {
		dbs := DBStructure{Chirps: map[int]Chirp{}}
		data, err := json.Marshal(dbs)
		if err != nil {
			return err
		}
		err = os.WriteFile(db.path, data, 0666)
		return err
	}
	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	dbs := DBStructure{}
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return DBStructure{}, err
	}
	return dbs, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0666)
	return err
}
