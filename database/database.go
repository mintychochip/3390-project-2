package database

import (
	"api-3390/config"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	_ "modernc.org/sqlite"
)

func getConnectionString(cfg *config.Config) string {
	return cfg.Path
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
	var connStr = getConnectionString(cfg)
	db, err := sql.Open(driverName, connStr)
	if err != nil {
		return nil, err
	}
	//TODO: handle this
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully connected to database")
	return db, nil
}
