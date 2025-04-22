package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GaussHammer/chirpy/internal/auth"
)

func userRefresh(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("No key in header!")
		w.WriteHeader(401)
		return
	}
	// Get user from valid token
	user, err := cfg.dbQueries.SelectUserFromToken(r.Context(), token)
	if err != nil {
		log.Printf("Invalid token: %v", err)
		w.WriteHeader(401)
		return
	}

	// Create new access token
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		w.WriteHeader(500)
		return
	}

	// Return new access token
	resp := struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
