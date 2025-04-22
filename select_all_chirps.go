package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func selectAllChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	dbChirps, err := cfg.dbQueries.SelectChirps(r.Context())
	if err != nil {
		log.Printf("error fetching the chirps: %s", err)
		w.WriteHeader(500)
		return
	}
	returnChirps := []returnVals{}
	for _, dbChirp := range dbChirps {
		returnChirps = append(returnChirps, returnVals{
			Id:        dbChirp.ID, // Assuming these field names match
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}
	dat, err := json.Marshal(returnChirps)
	if err != nil {
		log.Printf("couldn't marshal the JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
