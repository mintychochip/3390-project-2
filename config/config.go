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

//NewConfig
/*
Returns a new instance of the config, if the file path to a JSON configuration file is passed in the command-line,
the program will attempt to load from the file provided,
otherwise, it will attempt to load from environment variables.
*/
func NewConfig() (*Config, error) {
	if len(os.Args) > 1 {
		return load(os.Args[1])
	}
	return loadFromEnv()
}

//DatabaseConnection
/*
Returns a SQLite database connection using the parameters specified in the cfg `Config`.
*/
func (cfg *Config) DatabaseConnection() (*sql.DB, error) {
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

/*
Loads values from a JSON configuration file as specified by field tags in the `Config` struct.
*/
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

/*
loads configuration values from environment variables.

It retrieves the following variables:
- ADDRESS: The address to connect to.
- PATH: The path for the resource.
- REFERENCE_KEY: A key used for referencing.
- REFERENCE_HEADER: A header used for referencing.

Returns a pointer to a Config struct populated with these values,
or an error if any required environment variable is missing.
*/
func loadFromEnv() (*Config, error) {
	return &Config{
		Address:         os.Getenv("ADDRESS"),
		Path:            os.Getenv("PATH"),
		ReferenceKey:    os.Getenv("REFERENCE_KEY"),
		ReferenceHeader: os.Getenv("REFERENCE_HEADER"),
	}, nil
}
