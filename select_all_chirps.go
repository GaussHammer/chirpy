package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func selectAllChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	sortQuery := r.URL.Query().Get("sort")
	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		authorUuid := uuid.MustParse(authorID)
		authorChirps, err := cfg.dbQueries.SelectChirpsByUserId(r.Context(), authorUuid)
		if err != nil {
			log.Printf("Error fetching the author's chirps")
		}
		returnAuthorChirps := []returnVals{}
		for _, authorChirp := range authorChirps {
			returnAuthorChirps = append(returnAuthorChirps, returnVals{
				Id:        authorChirp.ID,
				CreatedAt: authorChirp.CreatedAt,
				UpdatedAt: authorChirp.UpdatedAt,
				Body:      authorChirp.Body,
				UserID:    authorChirp.UserID,
			})
		}
		dat, err := json.Marshal(returnAuthorChirps)
		if err != nil {
			log.Printf("error marshalling the chirps")
			w.WriteHeader(500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
		return
	}
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
	if sortQuery == "desc" {
		sort.Slice(returnChirps, func(i, j int) bool { return returnChirps[i].CreatedAt.After(returnChirps[j].CreatedAt) })
	} else if sortQuery == "asc" {
		sort.Slice(returnChirps, func(i, j int) bool { return returnChirps[i].CreatedAt.Before(returnChirps[j].CreatedAt) })
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
