package user

import (
	"database/sql"
	"fmt"
)

type Service[T any] interface {
	ExistsQuery(obj *T) (string, []interface{})
	ItemExists(obj *T) (bool, error)
}

func itemExists(db *sql.DB, query string, args []interface{}) (bool, error) {
	var exists bool
	err := db.QueryRow(query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error whilst checking existence: %w", err)
	}
	fmt.Println(exists)
	return exists, nil
}
