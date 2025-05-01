package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GaussHammer/chirpy/internal/auth"
	"github.com/google/uuid"
)

func polkaWebhooks(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type webHookBody struct {
		Event string `json:"event"`
		Data  struct {
			User_id uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := webHookBody{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding JSON")
		w.WriteHeader(500)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("error getting the API key")
	}

	if apiKey != cfg.polkaKey {
		w.WriteHeader(401)
		return
	}

	err = cfg.dbQueries.UpgradeChirpyRed(r.Context(), params.Data.User_id)
	if err != nil {
		log.Printf("couldn't upgrade the user")
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
}
