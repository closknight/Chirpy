package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/closknight/Chirpy/internal/auth"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := auth.HashPassword(userReq.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newUser, err := cfg.DB.CreateUser(userReq.Email, hash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, User{Id: newUser.Id, Email: newUser.Email})
}

func (cfg *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect login request")
	}
	user, err := cfg.DB.GetUserByEmail(userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := auth.CheckPasswordHash(user.HashPassword, userReq.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	ttl := 24 * time.Hour
	if userReq.ExpiresInSeconds > 0 && userReq.ExpiresInSeconds < int(ttl.Seconds()) {
		ttl = time.Duration(userReq.ExpiresInSeconds * int(time.Second))
	}

	tokenString, err := auth.CreateJWT(user.Id, cfg.jwtSecret, ttl)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email, Token: tokenString})
}

func (cfg apiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	tokenString, err := gettokenString(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	idString, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse token id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err = decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hash, err := auth.HashPassword(userReq.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.UpdateUser(id, userReq.Email, hash)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem updating user")
		return
	}
	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email})
}

func gettokenString(r *http.Request) (string, error) {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) < 2 || header[0] != "Bearer" {
		return "", errors.New("invalid header")
	}
	return header[1], nil
}
