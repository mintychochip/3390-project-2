package database

import (
	"api-3390/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	_ "modernc.org/sqlite"
)

func getConnectionString(cfg *config.Config, driverName string) string {
	if driverName == "sqlite" {
		return cfg.Path
	}
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.User, cfg.Password, cfg.DBName, cfg.Host, cfg.Port)
}

//	func InsertData(db *sql.DB) error {
//		query := 'INSERT INTO '
//		_, err := db.Exec(query, "Alice", 30)
//		if err != nil {
//			return err
//		}
//
// }
func Connection(cfg *config.Config) (*sql.DB, error) {
	driverName := "sqlite"
	var connStr = getConnectionString(cfg, driverName)
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, err
	}
	//TODO: handle this
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to database")
	return db, nil
}
