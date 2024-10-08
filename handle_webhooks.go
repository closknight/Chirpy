package main

import (
	"encoding/json"
	"net/http"

	"github.com/closknight/Chirpy/internal/auth"
)

func (cfg *apiConfig) HandlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}

	apiKey, err := auth.ParseAPIKEY(r.Header)
	if err != nil || apiKey != cfg.polka_api_key {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := parameters{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could't decode request")
		return
	}

	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.DB.UpgradeUsertoRed(request.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
