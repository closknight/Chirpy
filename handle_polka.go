package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) HandlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID int `json:"user_id"`
	}

	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	request := parameters{}
	err := decoder.Decode(&request)
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
