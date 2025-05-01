package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/GaussHammer/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
	polkaKey       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("failed to open connexion to database: %v", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is required")
	}
	apiCfg := apiConfig{
		dbQueries: database.New(db),
		platform:  os.Getenv("PLATFORM"),
		jwtSecret: jwtSecret,
		polkaKey:  polkaKey,
	}

	mux := http.NewServeMux()
	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		hits := apiCfg.fileserverHits.Load()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "text/html")
		w.Write([]byte(fmt.Sprintf(`<html>
			<body>
			  <h1>Welcome, Chirpy Admin</h1>
			  <p>Chirpy has been visited %d times!</p>
			</body>
		  </html>`, hits)))
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		if apiCfg.platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		apiCfg.fileserverHits.Store(0)
		apiCfg.dbQueries.DeleteUsers(r.Context())
	})
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		postChirp(w, r, &apiCfg)
	})
	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		selectAllChirps(w, r, &apiCfg)
	})
	mux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		selectChirp(w, r, &apiCfg)
	})
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {
		deleteChirp(w, r, &apiCfg)
	})
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		newUser(w, r, &apiCfg)
	})
	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		userLogin(w, r, &apiCfg)
	})
	mux.HandleFunc("POST /api/refresh", func(w http.ResponseWriter, r *http.Request) {
		userRefresh(w, r, &apiCfg)
	})
	mux.HandleFunc("POST /api/revoke", func(w http.ResponseWriter, r *http.Request) {
		userRevoke(w, r, &apiCfg)
	})
	mux.HandleFunc("POST /api/polka/webhooks", func(w http.ResponseWriter, r *http.Request) {
		polkaWebhooks(w, r, &apiCfg)
	})
	mux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, r *http.Request) {
		userUpdate(w, r, &apiCfg)
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
