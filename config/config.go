package config

import (
	"database/sql"
	"encoding/json"
	_ "modernc.org/sqlite"
	"os"
)

type Config struct {
	Address         string `json:"address"`
	Path            string `json:"path"`
	ReferenceKey    string `json:"reference_key"`
	ReferenceHeader string `json:"reference_header"`
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
	return &Config{
		Address:         os.Getenv("ADDRESS"),
		Path:            os.Getenv("PATH"),
		ReferenceKey:    os.Getenv("REFERENCE_KEY"),
		ReferenceHeader: os.Getenv("REFERENCE_HEADER"),
	}, nil
}
