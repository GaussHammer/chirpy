package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/GaussHammer/chirpy/internal/auth"
	"github.com/GaussHammer/chirpy/internal/database"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	resp := errorResponse{
		Error: msg,
	}
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func userLogin(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type connectUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := connectUser{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("couldn't decode JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	user, err := cfg.dbQueries.SelectUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return

	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(3600)*time.Second)
	if err != nil {
		respondWithError(w, 500, "Error creating JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "Error creating refresh token")
		return
	}
	expiresAt := sql.NullTime{
		Time:  time.Now().Add(60 * 24 * time.Hour),
		Valid: true,
	}

	cfg.dbQueries.CreateToken(r.Context(), database.CreateTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	})

	returnUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	dat, err := json.Marshal(returnUser)
	if err != nil {
		log.Printf("couldn't marshal into JSON: %s", err)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
