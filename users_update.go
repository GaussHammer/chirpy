package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GaussHammer/chirpy/internal/auth"
	"github.com/GaussHammer/chirpy/internal/database"
)

func userUpdate(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Wrong access code")
		w.WriteHeader(401)
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Invalid token")
		w.WriteHeader(401)
		return
	}
	type updateUser struct {
		NewEmail    string `json:"email"`
		NewPassword string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := updateUser{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding the body")
		w.WriteHeader(401)
		return
	}

	hashedPassword, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		log.Printf("error hashing new password")
		w.WriteHeader(401)
		return
	}
	userData, err := cfg.dbQueries.SelectUserById(r.Context(), userId)

	err = cfg.dbQueries.UpdateUserById(r.Context(), database.UpdateUserByIdParams{
		Email:          params.NewEmail,
		HashedPassword: hashedPassword,
		ID:             userId,
	})
	if err != nil {
		log.Printf("error updating the database")
		w.WriteHeader(401)
		return
	}
	returnUser := User{
		ID:           userId,
		CreatedAt:    userData.CreatedAt,
		UpdatedAt:    userData.UpdatedAt,
		Email:        params.NewEmail,
		AccessToken:  token,
		RefreshToken: userData.Token,
	}
	dat, err := json.Marshal(returnUser)
	if err != nil {
		log.Printf("couldn't marshal into JSON: %s", err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
