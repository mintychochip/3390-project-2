package main

import (
	"api-3390/auth"
	"api-3390/config"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

var userTable = `CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL
    );`

func main() {

	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := cfg.Connection()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(userTable)
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Post("/api/auth", func(w http.ResponseWriter, r *http.Request) {
		var u auth.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		b, err := auth.UserExists(db, &u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !b {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		authenticate, err := auth.UserAuthenticate(db, &u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !authenticate {
			http.Error(w, "user not authorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	r.Post("/api/register", func(w http.ResponseWriter, r *http.Request) {
		var u auth.User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			http.Error(w, fmt.Errorf("failed to register user: %w", err).Error(), http.StatusBadRequest)
			return
		}
		b, err := auth.UserExists(db, &u)
		if err != nil {
			http.Error(w, fmt.Errorf("failed to register user: %w", err).Error(), http.StatusBadRequest)
			return
		}
		if b {
			http.Error(w, "failed to register user: user already exists", http.StatusConflict)
			return
		}
		if err := auth.Register(db, u); err != nil {
			http.Error(w, fmt.Errorf("failed to register user: %w", err).Error(), http.StatusInternalServerError)
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
