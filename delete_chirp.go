package main

import (
	"log"
	"net/http"

	"github.com/GaussHammer/chirpy/internal/auth"
	"github.com/google/uuid"
)

func deleteChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	dbID := r.PathValue("chirpID")
	dbUUID, err := uuid.Parse(dbID)
	if err != nil {
		log.Printf("couldn't translate the string to UUID: %s", err)
		w.WriteHeader(500)
		return
	}
	dbChirp, err := cfg.dbQueries.SelectChirp(r.Context(), dbUUID)
	if err != nil {
		log.Printf("couldn't retrieve the chirp from the database")
		w.WriteHeader(404)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error authenticating user")
		w.WriteHeader(401)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("couldn't validate token")
		w.WriteHeader(401)
		return
	}
	if userId != dbChirp.UserID {
		log.Printf("Wrong user")
		w.WriteHeader(403)
		return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), dbChirp.ID)
	if err != nil {
		log.Printf("Error deleting chirp")
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
