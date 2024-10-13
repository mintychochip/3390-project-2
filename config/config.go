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
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

func LoadFromEnv() (Config, error) {
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	name := os.Getenv("NAME")

	if user == "" {
		return Config{}, errors.New("USER environment variable is not set")
	}
	if password == "" {
		return Config{}, errors.New("PASSWORD environment variable is not set")
	}
	if host == "" {
		return Config{}, errors.New("HOST environment variable is not set")
	}
	if name == "" {
		return Config{}, errors.New("NAME environment variable is not set")
	}

	return Config{
		User:     user,
		Password: password,
		Host:     host,
		Name:     name,
	}, nil
}
