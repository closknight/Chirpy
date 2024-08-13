package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(userID int, token string) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour * 60 * 24),
	}
	dbs.RefreshTokens[token] = refreshToken
	err = db.writeDB(dbs)
	return err
}

func (db *DB) GetUserFromRefreshToken(token string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	refreshToken, ok := dbs.RefreshTokens[token]
	if !ok {
		return User{}, errors.New("could not find token")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return User{}, errors.New("refresh token expired")
	}

	user, ok := dbs.Users[refreshToken.UserID]
	if !ok {
		return User{}, errors.New("could not find user")
	}
	return user, nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	dbs, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbs.RefreshTokens, token)
	err = db.writeDB(dbs)
	return err
}
