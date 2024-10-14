package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	User           string `json:"user"`
	Password       string `json:"password"`
	Host           string `json:"host"`
	DBName         string `json:"db_name"`
	Port           uint32 `json:"port"`
	SSLMode        bool   `json:"ssl_mode"`
	Path           string `json:"path"`
	ApplicationKey string `json:"application_key"`
}

func (cfg *Config) ApplicationMiddleWare(r *chi.Mux) {
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-KEY")
			if apiKey != cfg.ApplicationKey {
				http.Error(w, "Forbidden", http.StatusForbidden)
				clientIP := r.RemoteAddr
				if apiKey == "" {
					log.Printf("Unauthorized access attempt from IP: %s with empty key", clientIP)
					return
				}
				log.Printf("Unauthorized access attempt from IP: %s with key: %s\n", clientIP, apiKey)
				return
			}
			next.ServeHTTP(w, r)
		})
	})
}
func (cfg *Config) Connection() (*sql.DB, error) {
	driverName := "sqlite"
	var connStr = cfg.Path
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (cfg *Config) Address() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
func Load(path string) (*Config, error) {
	cfg, err := loadFromFile(path)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadFromFile(path string) (*Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadFromEnv() (*Config, error) {
	cfg := Config{
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Host:     os.Getenv("HOST"),
		DBName:   os.Getenv("DB_NAME"),
	}
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("'Port' cannot be empty:'%w'", err)
	}
	cfg.Port = uint32(port)
	return &cfg, nil
}
