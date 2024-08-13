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

	_, dupe := dbs.findUserByEmail(email)
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
	return user, nil
}

func (dbs *DBStructure) findUserByEmail(email string) (User, bool) {
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
	user, ok := dbs.findUserByEmail(email)
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbs.Users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil

}

func (db *DB) UpdateUser(userID int, email, hashedPassword string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbs.Users[userID]
	if !ok {
		return User{}, errors.New("user does not exist")
	}
	if dupe, ok := dbs.findUserByEmail(email); ok && dupe.Id != userID {
		return User{}, errors.New("email already in used")
	}

	user.Email = email
	user.HashPassword = hashedPassword
	dbs.Users[userID] = user

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
