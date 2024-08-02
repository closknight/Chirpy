package database

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, fs.ErrNotExist) {
		dbs := DBStructure{
			Chirps: map[int]Chirp{},
			Users:  map[int]User{},
		}
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
