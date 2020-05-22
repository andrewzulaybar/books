package postgres

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/lib/pq" // postgres driver
)

const path string = "internal/sql/"

var tables []string = []string{
	"location",
	"account_user",
	"author",
	"work",
	"publication",
}

// DB wraps our SQL database to allow for mocking.
type DB struct{ *sql.DB }

// Connect creates and returns a pool of connections to the database.
func Connect(params string) (*DB, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Disconnect closes the pool of connections to the given database.
func Disconnect(db *DB) {
	db.Close()
}

// Init creates tables by running the appropriate SQL scripts
// and also inserts existing data that we have into the tables.
func Init(db *DB) error {
	for _, path := range getFileNames() {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		statements := string(bytes)
		if _, err := db.Exec(statements); err != nil {
			return err
		}
	}
	return nil
}

func getFileNames() []string {
	fileNames := []string{path + "init.sql"}
	for _, table := range tables {
		fileName := path + table + ".sql"
		fileNames = append(fileNames, fileName)
	}
	return fileNames
}
