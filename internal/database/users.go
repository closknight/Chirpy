package database

import (
	"errors"
)

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	HashPassword string `json:"hash_password"`
}

func (db *DB) CreateUser(email, password string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	_, dupe := dbs.findUser(email)
	if dupe {
		return User{}, errors.New("username already exists")
	}

	id := len(dbs.Users) + 1
	user := User{
		Email:        email,
		Id:           id,
		HashPassword: password,
	}

	dbs.Users[id] = user
	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}
	resp := User{
		Email: user.Email,
		Id:    user.Id,
	}

	return resp, nil
}

func (dbs *DBStructure) findUser(email string) (User, bool) {
	for _, user := range dbs.Users {
		if user.Email == email {
			return user, true
		}
	}
	return User{}, false
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbs.findUser(email)
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}
