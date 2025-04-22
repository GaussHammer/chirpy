package main

import (
	"log"
	"net/http"
	"time"

	"github.com/GaussHammer/chirpy/internal/auth"
	"github.com/GaussHammer/chirpy/internal/database"
)

func userRevoke(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("No key in header!")
		w.WriteHeader(401)
		return
	}
	err = cfg.dbQueries.RevokeToken(r.Context(), database.RevokeTokenParams{
		UpdatedAt: time.Now(),
		Token:     token,
	})
	if err != nil {
		log.Printf("couldn't revoke the token")
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
