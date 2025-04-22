package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func selectChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
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
	returnChirp := returnVals{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	dat, err := json.Marshal(returnChirp)
	if err != nil {
		log.Printf("couldn't marshal into JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
