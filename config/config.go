package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Name     string `json:"name"`
}

func Load(path string) (Config, error) {
	cfg, err := loadFromFile(path)
	if err != nil {
		return Config{}, err
	}
	return cfg, validateConfig(&cfg)
}

func loadFromFile(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func LoadFromEnv() (Config, error) {
	cfg := Config{
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Host:     os.Getenv("HOST"),
		Name:     os.Getenv("NAME"),
	}
	return cfg, validateConfig(&cfg)
}

func validateConfig(cfg *Config) error {
	if cfg.User == "" {
		return errors.New("'User' is not set")
	}
	if cfg.Password == "" {
		return errors.New("'Password' is not set")
	}
	if cfg.Host == "" {
		return errors.New("'Host' is not set")
	}
	if cfg.Name == "" {
		return errors.New("'Name' is not set")
	}
	return nil
}
