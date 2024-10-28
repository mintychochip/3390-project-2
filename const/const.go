package constants

import (
	"api-3390/container/predicate"
	"io"
)

var FileMap = map[string][]predicate.Predicate[io.Reader]{
	".csv": {
		predicate.CSVContainsNumeric,
	},
}

const UserTable = `CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(128) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL
    );`
const UserFileTable = `CREATE TABLE IF NOT EXISTS user_files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    upload_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
)`
