package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"slices"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	database := DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := database.ensureDB()
	return &database, err
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

	dbs.Chirps[len(dbs.Chirps)+1] = chirp
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, fs.ErrNotExist) {
		dbs := DBStructure{Chirps: map[int]Chirp{}}
		return db.writeDB(dbs)
	}
	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbs := DBStructure{}

	data, err := os.ReadFile(db.path)
	if err != nil {
		return dbs, err
	}
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return dbs, err
	}
	return dbs, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, data, 0600)
	return err
}
