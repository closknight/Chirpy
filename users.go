package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/closknight/Chirpy/internal/auth"
	"github.com/golang-jwt/jwt/v5"
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

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "Chirpy",
		Subject:   strconv.Itoa(user.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email, Token: ss})
}

func (cfg apiConfig) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	tokenString := gettokenString(r)
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect token")
		return
	}
	idString, err := token.Claims.GetSubject()
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "users token has no id")
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not parse token id")
		return
	}
	user, err := cfg.DB.GetUserById(id)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not find users id")
		return
	}

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	userReq := parameters{}
	err = decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(userReq.Email) > 0 {
		user.Email = userReq.Email
	}
	if len(userReq.Password) > 0 {
		hash, err := auth.HashPassword(userReq.Password)

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		user.HashPassword = hash
	}

	user, err = cfg.DB.UpdateUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Problem updating user")
		return
	}
	respondWithJSON(w, http.StatusOK, User{Id: user.Id, Email: user.Email})
}

func gettokenString(r *http.Request) string {
	header := strings.Split(r.Header.Get("Authorization"), " ")
	if len(header) != 2 {
		return ""
	}
	return header[1]
}
