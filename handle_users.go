package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/closknight/Chirpy/internal/auth"
)

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
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
	respondWithJSON(
		w,
		http.StatusCreated,
		User{Id: newUser.Id, Email: newUser.Email, IsChirpyRed: newUser.IsChirpyRed})
}

func (cfg *apiConfig) HandleRefreshLogin(w http.ResponseWriter, r *http.Request) {
	type Resp struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.ParseBearer(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find token")
		return
	}
	user, err := cfg.DB.GetUserFromRefreshToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	access_token, err := auth.CreateJWT(user.Id, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create access token")
		return
	}
	respondWithJSON(w, http.StatusOK, Resp{Token: access_token})
}

func (cfg *apiConfig) HandleRevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.ParseBearer(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "could not find token")
		return
	}
	err = cfg.DB.RevokeRefreshToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not revoke token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
		return
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

	ttl := time.Hour
	if userReq.ExpiresInSeconds > 0 && userReq.ExpiresInSeconds < int(ttl.Seconds()) {
		ttl = time.Duration(userReq.ExpiresInSeconds * int(time.Second))
	}

	tokenString, err := auth.CreateJWT(user.Id, cfg.jwtSecret, ttl)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := auth.CreateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token")
		return
	}

	err = cfg.DB.SaveRefreshToken(user.Id, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not save refresh token")
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		User{
			Id:           user.Id,
			Email:        user.Email,
			Token:        tokenString,
			RefreshToken: refreshToken,
			IsChirpyRed:  user.IsChirpyRed,
		})
}

func (cfg apiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	tokenString, err := auth.ParseBearer(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	idString, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
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
	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email, IsChirpyRed: user.IsChirpyRed})
}
