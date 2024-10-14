package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	DBName   string `json:"db_name"`
	Port     uint32 `json:"port"`
	SSLMode  bool   `json:"ssl_mode"`
	Path     string `json:"path"`
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
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

//
//func validateConfig(cfg *Config) error {
//	if cfg.User == "" {
//		return errors.New("'User' is not set")
//	}
//	if cfg.Password == "" {
//		return errors.New("'Password' is not set")
//	}
//	if cfg.Host == "" {
//		return errors.New("'Host' is not set")
//	}
//	if cfg.DBName == "" {
//		return errors.New("'DBName' is not set")
//	}
//	return nil
//}
