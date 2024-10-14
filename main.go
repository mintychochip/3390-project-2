package main

import (
	"api-3390/auth"
	"api-3390/config"
	"api-3390/database"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {

	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Connection(cfg)
	r := chi.NewRouter()
	r.Post("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		type LoginRequest struct {
			Field      string `json:"field"`
			Identifier string `json:"identifier"`
			Password   string `json:"password"`
		}
		var request LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(request); err != nil {
			http.Error(w, "Error registering user", http.StatusInternalServerError)
			return
		}
		authenticate, err := auth.Authenticate(db, request.Field, request.Identifier, request.Password)
		if err != nil {
			http.Error(w, "Error authenticating user", http.StatusInternalServerError)
			return
		}
		if !authenticate {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/api/register", func(w http.ResponseWriter, r *http.Request) {
		var user auth.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err := auth.Register(db, user); err != nil {
			var str = fmt.Errorf("Error registering user: %w", err).Error()
			http.Error(w, str, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
	log.Println(fmt.Sprintf("Starting server on: '%s'", cfg.Address()))
	if err := http.ListenAndServe(cfg.Address(), r); err != nil {
		log.Fatal(err)
	}
}

func getConfig() (*config.Config, error) {
	if len(os.Args) > 1 {
		cfg, err := config.Load(os.Args[1])
		return cfg, err
	}
	cfg, err := config.LoadFromEnv()
	return cfg, err

}
