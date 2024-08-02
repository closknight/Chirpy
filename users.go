package main

import (
	"encoding/json"
	"net/http"

	"github.com/closknight/Chirpy/internal/auth"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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
		Password string `json:"password"`
		Email    string `json:"email"`
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

	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email})
}
