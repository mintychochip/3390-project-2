package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"strconv"
)

type Config struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Host         string `json:"host"`
	DBName       string `json:"db_name"`
	Port         uint32 `json:"port"`
	SSLMode      bool   `json:"ssl_mode"`
	Path         string `json:"path"`
	ReferenceKey string `json:"reference_key"`
}

func NewConfig() (*Config, error) {
	if len(os.Args) > 1 {
		return load(os.Args[1])
	}
	return loadFromEnv()
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
func load(path string) (*Config, error) {
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

func loadFromEnv() (*Config, error) {
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
