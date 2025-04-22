package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/GaussHammer/chirpy/internal/auth"
	"github.com/GaussHammer/chirpy/internal/database"
	"github.com/google/uuid"
)

type returnVals struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func postChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("couldn't get the token from header")
		w.WriteHeader(401)
		return
	}
	auth, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("invalid token")
		w.WriteHeader(401)
		return
	}

	var profane_words = []string{"kerfuffle", "sharbert", "fornax"}
	respBody := returnVals{}
	if len(params.Body) > 140 {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(400)
		return
	} else {
		var cleanedBody string
		splitBody := strings.Split(params.Body, " ")
		for i := 0; i < len(splitBody); i++ {
			for j := 0; j < len(profane_words); j++ {
				if strings.ToLower(splitBody[i]) == profane_words[j] {
					splitBody[i] = "****"
				}
			}
		}
		cleanedBody = strings.Join(splitBody, " ")

		chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: auth,
		})
		if err != nil {
			log.Printf("Error inserting chirp in database")
			w.WriteHeader(500)
			return
		}
		respBody = returnVals{
			Id:        chirp.ID,
			Body:      cleanedBody,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    auth,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(201)
		w.Write(dat)
	}
}
